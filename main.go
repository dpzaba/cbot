package main

import (
	"flag"
	"log"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/cabify/cbot/flowdock"
)

var (
	prefix      = flag.String("prefix", "cbot", "bot prefix")
	token       = flag.String("token", "", "rest token")
	flows       = flag.String("flows", "", "flows, separted by comma")
	commandsDir = flag.String("c", "commands", "commands directory")
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

func handleMessage(c *flowdock.Client, e flowdock.Event, responders []*MessageResponder) {
	content, args, err := parseMessageContent(e)
	if err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}
	if len(content) == 0 {
		return
	}
	direct := len(args) > 0 && args[0] == *prefix
	directHandled := !direct
	for _, responder := range responders {
		caught, err := responder.Handle(direct, content, args[1:], func(response string) error {
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
	if !directHandled {
		log.Printf("Unhandled direct message: %s", content)
	}
}

var spaceSplitter *regexp.Regexp = regexp.MustCompile("\\s+")

func parseMessageContent(e flowdock.Event) (string, []string, error) {
	content, err := e.MessageContent()
	if err != nil {
		return content, nil, err
	}
	cleaned := strings.ToLower(strings.TrimSpace(content))
	args := spaceSplitter.Split(cleaned, -1)
	return content, args, nil
}

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
