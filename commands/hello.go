package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var msg string
	if rand.Float32() > 0.9 {
		msg = "Leave me alone!"
	} else {
		msg = "Hello from the internetz!"
	}
	fmt.Println(msg)
}
