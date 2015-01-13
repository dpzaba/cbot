package main

import (
	"flag"
	"fmt"
	"github.com/franela/goreq"
	"log"
	"strconv"
	// "strings"
	"time"
)

const moustache_api string = "http://mustachify.me/?src="
const github_api_url string = "https://api.github.com/"

type GithubUser struct {
	Login      string
	Avatar_url string
}

func init() {
	goreq.SetConnectTimeout(10 * time.Second)
	flag.Parse()
}

func main() {
	args := flag.Args()

	if len(args) < 1 {
		help()
		return
	}

	cmd := args[0]

	switch cmd {
	case "gh":

		if len(args) < 2 {
			fmt.Println("Github user not provided OMG")
			return
		}

		gh := args[1]

		avatar_url := getGithubAvatar(gh)
		fmt.Println(moustache_api + avatar_url)
		return

	case "help":
		help()
		return
	}

	fmt.Println(moustache_api + args[0])

}

func help() {
	fmt.Println("Available commands")
	fmt.Println("mustachify gh {github_user_name}")
	fmt.Println("mustachify {img_url}")
}

func getGithubAvatar(user string) string {
	endpoint := github_api_url + "users/" + user
	res, err := goreq.Request{Uri: endpoint}.Do()
	panicErr(err)
	assertSucess(res.StatusCode)
	var userEncoder GithubUser
	res.Body.FromJsonTo(&userEncoder)
	return userEncoder.Avatar_url
}

func panicErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func assertSucess(code int) {
	if code != 200 {
		log.Fatal("Response code: " + strconv.Itoa(code))
	}
}
