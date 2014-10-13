package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"bufio"
)

func InitMessageResponders(dir string) ([]*MessageResponder, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(absDir)
	if err != nil {
		return nil, err
	}
	responders := make([]*MessageResponder, 0, len(files))
	for _, f := range files {
		if f.IsDir() || f.Mode()&0111 == 0 { // assert not a dir and is exectuable
			continue
		}
		path := filepath.Join(absDir, f.Name())
		name := strings.ToLower(f.Name())
		responders = append(responders, &MessageResponder{
			Name: name,
			Path: path,
		})
	}
	return responders, nil
}

type MessageResponder struct {
	Name string
	Path string
}

func (c *MessageResponder) Handle(direct bool, content string, args []string, responder func(string) error) (bool, error) {
	var cmd *exec.Cmd
	if direct && strings.HasPrefix(c.Name, "_") {
		cmd = exec.Command(c.Path, content)
	} else if len(args) >= 1 && args[0] == c.Name {
		cmd = exec.Command(c.Path, args[1:]...)
	} else {
		return false, nil
	}
	cmd.Dir = filepath.Dir(c.Path)
	output, err := runCommand(cmd)
	if err != nil {
		return true, err
	}
	go func() {
		for o := range output {
			if err := responder(o); err != nil {
				log.Printf("Error responding to command: %v", err)
				return
			}
		}
	}()
	return true, nil
}

// runs command, returning stdout and stderr on chan
func runCommand(cmd *exec.Cmd) (<-chan string, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	outScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	output := make(chan string, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for outScanner.Scan() {
			output <- outScanner.Text()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for errScanner.Scan() {
			output <- fmt.Sprintf("Err: %s", errScanner.Text())
		}
	}()
	go func() {
		defer close(output)
		if err := cmd.Wait(); err != nil {
			log.Printf("Error on cmd.Wait: %v", err)
		}
		wg.Wait()
	}()
	return output, nil
}

/*
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
				cmd := exec.Command(exePath, args...)
				cmd.Dir = absDir
				output, err := outputChan(cmd)
				if err != nil {
					outputHandler(e, fmt.Sprintf("Error running command: %s", err))
					return err
				}
				for o := range output {
					if len(o) > 0 {
						if err := outputHandler(e, o); err != nil {
							return err
						}
					}
				}
				return nil
			},
		}
		fmt.Printf("Registered <%s %s> command\n", prefix, name)
		responders = append(responders, r)
	}
	return responders, nil
}

// runs command, returning stdout and stderr on chan
func outputChan(cmd *exec.Cmd) (<-chan string, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	output := make(chan string, 1)
	outScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	go func() {
		for outScanner.Scan() {
			output <- outScanner.Text()
		}
	}()
	go func() {
		for errScanner.Scan() {
			output <- fmt.Sprintf("err: %s", errScanner.Text())
		}
	}()
	go func() {
		defer close(output)
		if err := cmd.Wait(); err != nil {
			log.Printf("Error on cmd.Wait: %v", err)
		}
	}()
	return output, nil
}

type MessageCommandResponder struct {
	Prefix string
	Name   string
	Run    func(event Event, content string, args []string) error
}

func (m *MessageCommandResponder) Handle(event Event, content string, args []string) {
	if len(args) >= 2 && args[0] == m.Prefix && args[1] == m.Name {
		if err := m.Run(event, content, args[2:]); err != nil {
			log.Println(err)
		}
	}
}
*/
