package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"regexp"
)

func init() {
	flag.Parse()
}

func main() {
	args := flag.Args()
	message := args[0]
	r, _ := regexp.Compile("meme:([a-z]+)")
	memes := r.FindAllString(message, -1)
	for i := 0; i < len(memes); i++ {
		fmt.Println(getMeme(memes[i]))
	}
}

func getMeme(meme string) string {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Print("Failed to connect to redis")
		return ""
	}
	s, err := redis.String(c.Do("GET", "cbot:"+meme))
	if err != nil {
		fmt.Print("Meme not found")
		return ""
	}
	return s
}
