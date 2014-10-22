package main

import (
	"fmt"
	"github.com/franela/goreq"
	"strconv"
	"strings"
)

const api_url string = "https://hacker-news.firebaseio.com/v0/"

type TopStories []int

type Item struct {
	By    string
	Id    int
	Type  string
	Url   string
	Title string
	Text  string
	Score int
}

func main() {
	top := topstories()
	for i := 0; i < 10; i++ {
		id := top[i]
		printItem(item(id))
	}
}

func topstories() TopStories {
	var endpoint = "topstories.json"
	request_uri := api_url + endpoint
	res, err := goreq.Request{Uri: request_uri}.Do()
	panicErr(err)
	assertSucess(res.StatusCode)
	var storiesEncoder TopStories
	res.Body.FromJsonTo(&storiesEncoder)
	return storiesEncoder
}

func item(Id int) Item {
	id := strconv.Itoa(Id)
	var target = "item/#{id}.json"
	endpoint := strings.Replace(target, "#{id}", id, -1)
	request_uri := api_url + endpoint
	res, err := goreq.Request{Uri: request_uri}.Do()
	panicErr(err)
	assertSucess(res.StatusCode)
	var item Item
	res.Body.FromJsonTo(&item)
	return item
}

func printItem(Item Item) {
	var template = "(#{score})[#{by}]#{title} #{url}"
	r := strings.NewReplacer("#{score}", strconv.Itoa(Item.Score), "#{by}", Item.By, "#{title}", Item.Title, "#{url}", Item.Url)
	p := r.Replace(template)
	fmt.Println(p)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func assertSucess(code int) {
	if code != 200 {
		panic("Request failure")
	}
}
