package main

import "bitbucket.org/cabify/cbot/flowdock"

type MessageResponder interface {
	Handles(flowdock.Event, []string) bool
	Handle(flowdock.Event, []string)
}
