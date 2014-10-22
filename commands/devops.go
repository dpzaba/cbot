package main

import (
	"fmt"
	"github.com/franela/goreq"
	"encoding/json"
	"strings"
)

const tumblr_url_post = "http://api.tumblr.com/v2/blog/devopsreactions.tumblr.com/posts?id=#{id}&api_key=fuiKNFp9vQFvjLNvx4sUwti4Yb5yGutBN4Xh10LXZhhRKjWlV4"

type Post struct {
	Title string
	Url string `json:"short_url"`
}

type TumblrPostsResponse struct {
	Response struct {
		Blog json.RawMessage
		Posts []Post
	}
}

func main() {
	res, err := goreq.Request{ Uri: "http://devopsreactions.tumblr.com/random" }.Do()
	if err != nil {
		fmt.Println("Damn it yo, got an err")
		return
	}

	url := res.Header.Get("Location")
	id := strings.Split(url, "/")[4]

	if len(id) > 0 {
		request_url := strings.Replace(tumblr_url_post, "#{id}", id, -1)
		res, err := goreq.Request{ Uri: request_url }.Do()
		if err != nil {
			fmt.Println("D'oh, error response")
			return
		}
		var posts TumblrPostsResponse
		res.Body.FromJsonTo(&posts)
		fmt.Println(posts.Response.Posts[0].Title)
		fmt.Println(posts.Response.Posts[0].Url)
	}

}