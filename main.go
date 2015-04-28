package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/nlopes/slack"
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

	c := slack.New(*token)

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
func handleMessage(c *slack.Slack, e *slack.MessageEvent, responders []*MessageResponder) {
        if e.Username == "cbot" {
		return
	}

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
		params := slack.PostMessageParameters{ Username: *prefix }
		if _,_,err := c.PostMessage(e.ChannelId, helpTxt.String(), params); err != nil {
			log.Println(err)
		}
		return
	}

	directHandled := !direct

	user, err := c.GetUserInfo(e.Msg.UserId)

	os.Setenv("CURRENT_FLOW", e.ChannelId)
	os.Setenv("CURRENT_USER_AVATAR", user.Profile.ImageOriginal)
	os.Setenv("CURRENT_USER_EMAIL", user.Profile.Email)
	os.Setenv("CURRENT_USER_NICK", user.Name)
	os.Setenv("CURRENT_USER_NAME", user.Profile.RealName)
	os.Setenv("CURRENT_USER_ID", string(user.Id))

	for _, responder := range responders {

		caught, err := responder.Handle(direct, content, args[1:], func(response string) error {
			// handle the output of the command by replying to the message
			params := slack.PostMessageParameters{ Username: *prefix }
			_,_,error := c.PostMessage(e.ChannelId, response, params) 
			return error
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
		params := slack.PostMessageParameters{ Username: *prefix }
		if _,_,err := c.PostMessage(e.ChannelId, resp, params); err != nil {
			log.Println(err)
		}
	}
}

var spaceSplitter *regexp.Regexp = regexp.MustCompile("\\s+")

// parseMessageContent cleans a message's content and breaks into args
func parseMessageContent(e *slack.MessageEvent) (string, []string, error) {
	content := e.Msg.Text

	cleaned := strings.TrimSpace(content)
	args := spaceSplitter.Split(cleaned, -1)
	return content, args, nil
}

func startStream(c *slack.Slack, flows []string, responders []*MessageResponder) {
	chSender := make(chan slack.OutgoingMessage)
	chReceiver := make(chan slack.SlackEvent)

	wsAPI, err := c.StartRTM("", "https://slack.com/api")
	if err != nil {
		fmt.Errorf("%s\n", err)
	}
	go wsAPI.HandleIncomingEvents(chReceiver)
	go wsAPI.Keepalive(20 * time.Second)
	go func(wsAPI *slack.SlackWS, chSender chan slack.OutgoingMessage) {
		for {
			select {
			case msg := <-chSender:
				wsAPI.SendMessage(&msg)
			}
		}
	}(wsAPI, chSender)
	for {
		select {
		case msg := <-chReceiver:
			fmt.Print("Event Received: ")
			switch msg.Data.(type) {
			case slack.HelloEvent:
				// Ignore hello
			case *slack.MessageEvent:
				a := msg.Data.(*slack.MessageEvent)
				handleMessage(c, a, responders)
				fmt.Printf("Message: %v\n", a)
			case *slack.PresenceChangeEvent:
				a := msg.Data.(*slack.PresenceChangeEvent)
				fmt.Printf("Presence Change: %v\n", a)
			case slack.LatencyReport:
				a := msg.Data.(slack.LatencyReport)
				fmt.Printf("Current latency: %v\n", a.Value)
			case *slack.SlackWSError:
				error := msg.Data.(*slack.SlackWSError)
				fmt.Printf("Error: %d - %s\n", error.Code, error.Msg)
			default:
				fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}
