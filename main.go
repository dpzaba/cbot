package main

import (
	"encoding/json"
	"flag"
	"log"
	"regexp"
	"strings"
	"time"
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
	c, err := NewClient(*token)
	if err != nil {
		log.Fatal(err)
	}
	responders, err := InitExecutableCommands(*commandsDir, *prefix, func(e Event, output string) error {
		return c.Message(Message{
			Event:    "comment",
			Content:  output,
			Flow:     e.Flow,
			UserName: *prefix,
			ID:       e.ID,
		})
	})
	if err != nil {
		log.Fatal(err)
	}
	startStream(c, responders.handleEvent, strings.Split(*flows, ",")...)
}

type Responders []Responder

type Responder interface {
	Handle(event Event, content string, args []string) error
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
		if err := responder.Handle(event, content, args); err != nil {
			return err
		}
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
		return r.handleMessage(event)
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

func startStream(client *Client, handler func(msg []byte) error, flows ...string) {
	log.Printf("Connecting to %v\n", flows)
	stream, errors, err := client.Stream(flows...)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case msg, ok := <-stream:
			if !ok {
				log.Println("Stream closed!")
				time.Sleep(2 * time.Second)
				go startStream(client, handler, flows...)
				return
			}
			if err = handler(msg); err != nil {
				log.Printf("Handler error: %v", err)
			}
			//fmt.Println(string(msg))
		case err := <-errors:
			log.Printf("Stream Error: %v", err)
		}
	}
}
