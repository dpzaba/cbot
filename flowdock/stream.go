package flowdock

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Stream receives an array of flows (ie cabify/test, cabify/test2) and returns 2 channels
// a channel of raw messages and a channel of errors produced while parsing the stream
func (c *Client) Stream(flows []string) (<-chan []byte, <-chan error, error) {
	endpoint := newStreamURL(c.token, flows)
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
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("stream response: %s. %s", resp.Status, string(body))
	}
	lines, errors := streamReq(resp.Body)
	return lines, errors, nil
}

// given an http.Response.Body, provide 2 channels:
//  - bytes per line received
//  - any errors while receiving stream
func streamReq(body io.ReadCloser) (<-chan []byte, <-chan error) {
	reader := bufio.NewReader(body)
	// channel which contains a slice of bytes per line received
	lines := make(chan []byte)
	// errors channel, with buffer so we don't block in case no readers
	errors := make(chan error, 1)
	blankLine := []byte{'\n'}
	// pump lines and errors into channel
	go func() {
		defer body.Close()
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				errors <- err
				close(lines)
				return
			}
			trimmed := bytes.TrimSpace(line)
			if len(trimmed) == 0 || bytes.Equal(blankLine, trimmed) {
				continue
			}
			lines <- trimmed
		}
	}()
	return lines, errors
}

func newStreamURL(token string, flows []string) string {
	base, _ := url.Parse(streamURL)
	base.User = url.User(token) // set token as BasicAuth user
	base.Path = "flows"
	base.RawQuery = (url.Values{
		"filter": []string{strings.Join(flows, ",")},
		"user":   []string{"0"}, // do not receive private messages to user
	}).Encode()
	return base.String()
}
