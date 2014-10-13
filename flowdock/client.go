package flowdock

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	token string // user's rest token
}

func NewClient(token string) (*Client, error) {
	return &Client{
		token: token,
	}, nil
}

/*
func (c *Client) StreamURL(flows ...string) string {
	base, _ := url.Parse(streamURL)
	base.User = url.User(c.token) // set token as BasicAuth user
	base.Path = "flows"
	params := url.Values{}
	params.Set("filter", strings.Join(flows, ","))
	params.Set("user", "0") // do not receive private messages
	base.RawQuery = params.Encode()
	return base.String()
}

func (c *Client) Stream(flows ...string) (<-chan []byte, <-chan error, error) {
	endpoint := c.StreamURL(flows...)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("stream response: %s", resp.Status)
	}
	reader := bufio.NewReader(resp.Body)
	lines := make(chan []byte)
	errors := make(chan error, 1)
	blankLine := []byte{'\n'}
	go func(r *bufio.Reader, c chan<- []byte, e chan<- error) {
		defer resp.Body.Close()
		for {
			line, err := r.ReadBytes('\n')
			if err != nil {
				e <- err
				close(c)
				return
			}
			if len(line) == 0 || bytes.Equal(blankLine, line) {
				continue
			}
			c <- line
		}
	}(reader, lines, errors)
	return lines, errors, nil
}

type Message struct {
	Event    string   `json:"event"`
	Content  string   `json:"content"`
	Flow     string   `json:"flow"`
	UserName string   `json:"external_user_name"`
	ID       int      `json:"message_id,omitempty"` // when responding to a message
	Tags     []string `json:"tags,omitempty"`
}

func (c *Client) Message(msg Message) error {
	endpoint, _ := url.Parse(restURL)
	endpoint.User = url.User(c.token) // set token as BasicAuth user
	if msg.Event == "comment" {
		endpoint.Path = "comments"
	} else {
		endpoint.Path = "messages"
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	//fmt.Println(endpoint.String(), string(data))
	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error sending chat message: %s. %s", resp.Status, string(body))
	}
	return nil
}
*/

func (c *Client) PostJSON(endpoint string, data []byte) error {
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error POSTing to Rest API: %s. %s", resp.Status, string(body))
	}
	return nil
}
