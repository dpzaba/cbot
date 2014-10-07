package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func InitExecutableCommands(dir string, prefix string, outputHandler func(e Event, output string) error) (Responders, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(absDir)
	if err != nil {
		return nil, err
	}
	responders := Responders{}
	for _, f := range files {
		if f.IsDir() || f.Mode()&0111 == 0 { // check if exectuable
			continue
		}
		exePath := filepath.Join(absDir, f.Name())
		name := strings.ToLower(f.Name())
		r := &MessageCommandResponder{
			Prefix: prefix,
			Name:   name,
			Run: func(e Event, content string, args []string) error {
				if err := outputHandler(e, fmt.Sprintf("Starting %s %v", name, args)); err != nil {
					return err
				}
				cmd := exec.Command(exePath, args...)
				output, err := cmd.CombinedOutput()
				if err != nil {
					outputHandler(e, fmt.Sprintf("Error: %s %s", err, output))
					return fmt.Errorf("%s %v %s. %s", exePath, args, err, output)
				}
				log.Printf("<%s> %s %v: %s", e.User, name, args, string(output))
				return outputHandler(e, string(output))
			},
		}
		fmt.Printf("Registered <%s %s> command\n", prefix, name)
		responders = append(responders, r)
	}
	return responders, nil
}

type MessageCommandResponder struct {
	Prefix string
	Name   string
	Run    func(event Event, content string, args []string) error
}

func (m *MessageCommandResponder) Handle(event Event, content string, args []string) error {
	if len(args) >= 2 && args[0] == m.Prefix && args[1] == m.Name {
		return m.Run(event, content, args[2:])
	}
	return nil
}
