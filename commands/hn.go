package main

import (
	"flag"
	"fmt"
	"github.com/franela/goreq"
	"log"
	"strconv"
	"strings"
	"time"
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

func init() {
	goreq.SetConnectTimeout(10 * time.Second)
	flag.Parse()
}

func main() {
	args := flag.Args()
	cmd := args[0]
	number, err := strconv.Atoi(args[1])
	if err != nil {
		number = 10
	}
	switch cmd {
	case "top":
		listTop(number)
		return
	}

	fmt.Println("Sorry, idk this command")
}

func listTop(Number int) {
	fmt.Println("Hacker News Top " + strconv.Itoa(Number))
	fmt.Println("=====================")
	top := topstories()
	for i := 0; i < Number; i++ {
		id := top[i]
		if id > 0 {
			fmt.Println(strconv.Itoa(i+1) + ":  " + item(id))
		}
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

func item(Id int) string {
	id := strconv.Itoa(Id)
	var target = "item/#{id}.json"
	endpoint := strings.Replace(target, "#{id}", id, -1)
	request_uri := api_url + endpoint
	res, err := goreq.Request{Uri: request_uri}.Do()
	if err != nil {
		return item(Id)
	}
	assertSucess(res.StatusCode)
	var item Item
	res.Body.FromJsonTo(&item)
	return parseItem(item)
}

func parseItem(Item Item) string {
	var template = "(#{score})[#{by}]#{title} #{url}"
	r := strings.NewReplacer("#{score}", strconv.Itoa(Item.Score), "#{by}", Item.By, "#{title}", Item.Title, "#{url}", Item.Url)
	p := r.Replace(template)
	return p
}

func panicErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func assertSucess(code int) {
	if code != 200 {
		log.Fatal("Response code: " + strconv.Itoa(code))
	}
}
