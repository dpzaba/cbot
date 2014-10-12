package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/cabify/cbot/flowdock"
)

var (
	token = flag.String("token", "", "rest token")
)

func main() {
	flag.Parse()
	if *token == "" {
		flag.PrintDefaults()
		return
	}
	c, err := flowdock.NewClient(*token)
	if err != nil {
		printError(err)
		os.Exit(1)
	}
	msg := flowdock.NewComment(11590, "0adc400f-ca1c-434b-81ee-c932f4fba2dd", "test", "hello")
	err = c.PostEvent(*msg)
	if err != nil {
		printError(err)
	}
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v", err)
}
