package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func init() {
	flag.Parse()
}

func main() {

	args := flag.Args()

	if len(args) < 2 {
		help()
		return
	}

	cmd := args[0]
	meme := args[1]

	switch cmd {

	case "set":
		memeText := args[2]
		set(meme, memeText)
		fmt.Print("Set meme " + meme)
		return
	case "del":
		del(meme)
		fmt.Print("Deleted meme " + meme)
		return
	case "get":
		get(meme)
		return
	}

	help()

}

func help() {
	fmt.Print("Meme command usage")
	fmt.Print("==================")
	fmt.Print("cbot meme [get|del] (meme-key)")
	fmt.Print("cbot meme [set] (meme-key) (meme-text|url)")
}

func set(meme string, text string) {
	c := conn()
	c.Do("SET", memeKey(meme), text)
}

func get(meme string) {
	c := conn()
	s, err := redis.String(c.Do("GET", memeKey(meme)))
	if err != nil {
		fmt.Print("Meme not found")
	} else {
		fmt.Print(s)
	}
}

func del(meme string) {
	c := conn()
	c.Do("DEL", memeKey(meme))
}

func conn() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	return c
}

func memeKey(meme string) string {
	return "meme:" + meme
}
