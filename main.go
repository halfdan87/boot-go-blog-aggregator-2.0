package main

import (
	"fmt"
	"os"

	"github.com/halfdan87/boot-go-blog-aggregator-2.0/internal/config"
)

func main() {
	fmt.Println("Starting...")

	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	s := &state{cfg: cfg}
	c := &commands{mapping: make(map[string]func(*state, command) error)}
	c.register("login", login)

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No command given")
		os.Exit(1)
	}
	if len(args) == 1 {
		fmt.Println("No username given")
		os.Exit(1)
	}
	c.run(s, command{name: args[0], args: args[1:]})

	/*
		for {
			fmt.Print("Command: ")
			cmd, err := readCommand()
			if err != nil {
				panic(err)
			}
			err = c.run(s, cmd)
			if err != nil {
				fmt.Println(err)
			}
		}
	*/
}

func login(s *state, cmd command) error {
	return s.cfg.SetUser(cmd.args[0])
}

func readCommand() (command, error) {
	var cmd command
	var err error
	cmd.name, err = readString()
	if err != nil {
		return cmd, err
	}
	cmd.args, err = readArgs()
	return cmd, err
}

func readString() (string, error) {
	var s string
	_, err := fmt.Scanln(&s)
	return s, err
}

func readArgs() ([]string, error) {
	var args []string
	for {
		var arg string
		_, err := fmt.Scanln(&arg)
		if err != nil {
			return args, err
		}
		args = append(args, arg)
	}
}

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	mapping map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if _, ok := c.mapping[name]; ok {
		return fmt.Errorf("command %v already registered", name)
	}
	c.mapping[name] = f
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.mapping[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command %v", cmd.name)
	}
	return f(s, cmd)
}
