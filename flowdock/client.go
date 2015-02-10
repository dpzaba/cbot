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
		return fmt.Errorf("error POSTing '%s' to Rest API: %s. %s", string(data), resp.Status, string(body))
	}
	return nil
}

func (c *Client) GetJSON(endpoint string) (string, error) {
	req, err := http.NewRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("error Getting '%s' to Rest API: %s. %s", string(endpoint), resp.Status, string(body))
	} else {
		return string(body)
	}
}
