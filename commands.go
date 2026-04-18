package main

import (
	"errors"
)

type command struct {
	name 		string
	arguments	[]string
}

type commands struct {
	list map[string]func(*state, command) error
}

func (c commands) run(s *state, cmd command) error {
	if function, ok := c.list[cmd.name]; ok {
		return function(s, cmd)
	}
	return errors.New("command not supported")
}

func (c commands) register(name string, f func(*state, command) error) error {
	if _, exists := c.list[name]; exists {
		return errors.New("command already exists")
	}
	
	c.list[name] = f

	return nil
}