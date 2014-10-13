package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/cabify/cbot/flowdock"
)

var (
	token = flag.String("token", "", "rest token")
	flows = flag.String("flows", "", "flows, separted by comma")
)

func main() {
	flag.Parse()
	if *token == "" || *flows == "" {
		flag.PrintDefaults()
		return
	}
	c, err := flowdock.NewClient(*token)
	if err != nil {
		printError(err)
		os.Exit(1)
	}
	raw, errors, err := c.Stream(strings.Split(*flows, ","))
	if err != nil {
		printError(err)
		os.Exit(1)
	}
	for {
		select {
		case line, ok := <-raw:
			if !ok {
				fmt.Println("stream closed")
				return
			}
			fmt.Println(string(line))
		case err := <-errors:
			printError(err)
		}
	}
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v", err)
}
