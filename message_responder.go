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

// InitMessageResponders scans given directory and creates MessageResponders from the exectuable files
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

// Handle receives a string message and decides whether to respond to the message
func (c *MessageResponder) Handle(direct bool, content string, args []string, responder func(string) error) (bool, error) {
	// if the executable file begins with _, that can respond to any message and is therefore passed every message
	// otherwise, the filename must match the command
	var cmd *exec.Cmd
	if direct && len(args) > 0 && args[0] == c.Name {
		cmd = exec.Command(c.Path, args[1:]...)
	} else if !direct && c.Name[0] == '_' {
		cmd = exec.Command(c.Path, content)
	} else {
		return false, nil
	}
	// set the exe's CWD to it's directory
	cmd.Dir = filepath.Dir(c.Path)
	output, err := runCommand(cmd)
	if err != nil {
		return true, err
	}
	go func() {
		// listen to the output generated and pass it to the responder func
		for o := range output {
			if len(o) == 0 { // ignore blank lines
				continue
			}
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
	outScanner.Split(scanDoubleLines)
	errScanner := bufio.NewScanner(stderr)
	errScanner.Split(scanDoubleLines)
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	output := make(chan string, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for outScanner.Scan() {
			output <- strings.TrimRightFunc(outScanner.Text(), unicode.IsSpace)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for errScanner.Scan() {
			output <- fmt.Sprintf("STDERR: %s", errScanner.Text())
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

func scanDoubleLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := strings.Index(string(data), "\n\n"); i >= 0 {
		// We have a full newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}
