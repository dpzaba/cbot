package main

import (
	"flag"
	"log"
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
	startStream(c, rooms)
}

/*
type Responders []Responder

type Responder interface {
	Handle(event Event, content string, args []string)
}

var spaceSplitter *regexp.Regexp = regexp.MustCompile("\\s+")

func (r *Responders) handleMessage(event Event) error {
	var content string
	if err := json.Unmarshal(event.Content, &content); err != nil {
		return err
	}
	cleaned := strings.ToLower(strings.TrimSpace(content))
	args := spaceSplitter.Split(cleaned, -1)
	for _, responder := range *r {
		go responder.Handle(event, content, args)
	}
	return nil
}

func (r *Responders) handleEvent(msg []byte) error {
	var event Event
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}
	switch event.Event {
	case "message": // normal message (not threaded)
		if err := r.handleMessage(event); err != nil {
			return err
		}
	}
	return nil
}

type Event struct {
	Event   string          `json:"event"`
	ID      int             `json:"id,omitempty"`
	User    string          `json:"user,omitempty"`
	Flow    string          `json:"flow,omitempty"`
	Content json.RawMessage `json:"content"`
}

*/

func startStream(c *flowdock.Client, flows []string) {
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
				go startStream(c, flows)
				return
			}
			switch event.Event {
			case "message":
				handleMessage(event)
			}
		case err := <-errors:
			log.Printf("Stream Error: %v", err)
		}
	}
}
