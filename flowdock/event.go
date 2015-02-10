package flowdock

import (
	"encoding/json"
	"net/url"
)

func (c *Client) EventStream(flows []string) (<-chan Event, <-chan error, error) {
	data, streamErrs, err := c.Stream(flows)
	if err != nil {
		return nil, nil, err
	}
	errs := make(chan error, 1)
	go func() {
		for e := range streamErrs {
			errs <- e
		}
	}()
	events := make(chan Event)
	go func() {
		for l := range data {
			var event Event
			if err := json.Unmarshal(l, &event); err != nil {
				errs <- err
				continue
			}
			events <- event
		}
		close(events)
	}()
	return events, errs, nil
}

func NewMessage(flow, username, content string) *Event {
	return (&Event{
		Event:    "message",
		UserName: username,
		Flow:     flow,
	}).StringContent(content)
}

func NewComment(message uint, flow, username, content string) *Event {
	return (&Event{
		Event:    "comment",
		UserName: username,
		Flow:     flow,
		Message:  message,
	}).StringContent(content)
}

// SendMessage posts a message
func (c *Client) PostEvent(e Event) error {
	endpoint, _ := url.Parse(restURL)
	endpoint.User = url.User(c.token) // set token as BasicAuth user
	// depending on if this is a reply or message, the endpoint varies
	if e.Event == "comment" {
		endpoint.Path = "comments"
	} else {
		endpoint.Path = "messages"
	}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return c.PostJSON(endpoint.String(), data)
}

func (c *Client) GetUserById(id uint) (User, error) {
	endpoint, _ := url.Parse(restURL)
	endpoint.User = url.User(c.token) // set token as BasicAuth user
	// depending on if this is a reply or message, the endpoint varies
	endpoint.Path = "user/" + string(id)

	user := &User{}
	body, error := c.GetJSON(endpoint.String())
	json.Unmarshal(body, &user)

	return user, nil
}

// {"event":"message","tags":[],"uuid":"GfcU296O-g6rdUU4","id":11295,"flow":"0adc400f-ca1c-434b-81ee-c932f4fba2dd","content":"cool no hurries","sent":1412938667994,"app":"chat","attachments":[],"user":"106208"}
// {"event":"activity.user","tags":[],"uuid":null,"persist":false,"id":11301,"flow":"0adc400f-ca1c-434b-81ee-c932f4fba2dd","content":{"last_activity":1412938720917},"sent":1412938741475,"app":null,"attachments":[],"user":"104062"}

type Event struct {
	Event    string           `json:"event"`
	Content  *json.RawMessage `json:"content"`
	Flow     string           `json:"flow"`
	UserName string           `json:"external_user_name,omitempty"` // send when sending as a different user, ie "cbot"
	User     string           `json:"user,omitempty"`
	ID       uint             `json:"id,omitempty"`
	Message  uint             `json:"message,omitempty"` // required for comments, The id of the commented parent message (which must not be a comment)
	UUID     string           `json:"uuid,omitempty"`
	Tags     []string         `json:"tags,omitempty"`
	//	Sent     uint             `json:"sent,omitempty"`
}

type User struct {
	Email   string `json:email`
	Name    string `json:name`
	Avatar  string `json:avatar`
	Id      uint   `json:id`
	Website string `json:website`
	Nick    string `json:nick`
}

// StringContent sets the content to a string
func (e *Event) StringContent(content string) *Event {
	raw, err := json.Marshal(content)
	if err != nil {
		return e
	}
	r := json.RawMessage(raw)
	e.Content = &r
	return e
}

// MessageContent parse an Event.Content as a string
func (e Event) MessageContent() (string, error) {
	var c string
	err := json.Unmarshal(*e.Content, &c)
	return c, err
}
