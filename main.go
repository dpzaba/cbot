package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"bytes"

	"bitbucket.org/cabify/cbot/flowdock"
)

var (
	prefix      = flag.String("prefix", "cbot", "bot prefix")          // prefix for direct commands
	token       = flag.String("token", "", "rest token")               // Flowdock rest token
	flows       = flag.String("flows", "", "flows, separted by comma") // of form cabify/test,cabify/other
	commandsDir = flag.String("c", "commands", "commands directory")   // directory of the executable commands
)

func main() {
	flag.Parse()
	if *token == "" || *flows == "" {
		flag.PrintDefaults()
		return
	}
	c, err := flowdock.NewClient(*token)
	if err != nil {
		log.Fatal(err)
	}
	rooms := strings.Split(*flows, ",")
	responders, err := InitMessageResponders(*commandsDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, responder := range responders {
		log.Printf("Registered <%s> responder", responder.Name)
	}
	startStream(c, rooms, responders)
}

// handleMessage receives an event and passes it to the MessageResponders
func handleMessage(c *flowdock.Client, e flowdock.Event, responders []*MessageResponder) {
	content, args, err := parseMessageContent(e)
	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}
	if len(content) == 0 {
		return
	}
	// determine if this is a direct message (prefixed with bots name)
	direct := len(args) > 0 && args[0] == *prefix

	// handle 'help' command
	if direct && (len(args) == 1 || (len(args) > 1 && args[1] == "help")) {
		helpTxt := bytes.NewBufferString("I understand:\n")
		for _, r := range responders {
			if r.Name[0] == '_' {
				continue
			}
			helpTxt.WriteString(fmt.Sprintf("    %s %s\n", *prefix, r.Name))
		}
		comment := flowdock.NewComment(e.ID, e.Flow, *prefix, helpTxt.String())
		if err := c.PostEvent(*comment); err != nil {
			log.Println(err)
		}
		return
	}

	directHandled := !direct
	for _, responder := range responders {
		caught, err := responder.Handle(direct, content, args[1:], func(response string) error {
			// handle the output of the command by replying to the message
			comment := flowdock.NewComment(e.ID, e.Flow, *prefix, response)
			return c.PostEvent(*comment)
		})
		if err != nil {
			log.Println(err)
			continue
		}
		if caught && direct {
			directHandled = true
		}
	}
	// handle case when a direct message wasn't handled
	if !directHandled {
		log.Printf("Unhandled direct message: %s", content)
		resp := "Sorry, didn't recognize that command. Try 'cbot help'"
		comment := flowdock.NewComment(e.ID, e.Flow, *prefix, resp)
		if err := c.PostEvent(*comment); err != nil {
			log.Println(err)
		}
	}
}

var spaceSplitter *regexp.Regexp = regexp.MustCompile("\\s+")

// parseMessageContent cleans a message's content and breaks into args
func parseMessageContent(e flowdock.Event) (string, []string, error) {
	content, err := e.MessageContent()
	if err != nil {
		return content, nil, err
	}
	cleaned := strings.TrimSpace(content)
	args := spaceSplitter.Split(cleaned, -1)
	return content, args, nil
}

// startStream starts streaming the given flows and responds to messages
func startStream(c *flowdock.Client, flows []string, responders []*MessageResponder) {
	log.Printf("Connecting to %v\n", flows)
	events, errors, err := c.EventStream(flows)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case event, ok := <-events:
			if !ok {
				log.Println("Stream closed! Restarting")
				time.Sleep(2 * time.Second)
				go startStream(c, flows, responders)
				return
			}
			switch event.Event {
			case "message":
				handleMessage(c, event, responders)
			}
		case err := <-errors:
			log.Printf("Stream Error: %v", err)
		}
	}
}
