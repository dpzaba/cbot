package main

import (
	"os/exec"

	"bitbucket.org/cabify/cbot/flowdock"
)

type CommandMessageResponder struct {
	Name string
	Cmd  *exec.Cmd
}

func (c *CommandMessageResponder) Handles(e flowdock.Event, args []string) bool {
	return len(args) > 0 && args[0] == c.Name
}

func (c *CommandMessageResponder) Handle(e flowdock.Event, args []string) {

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
