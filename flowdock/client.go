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
