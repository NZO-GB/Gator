package main

import (
	"fmt"
	"errors"
	config "github.com/NZO-GB/Gator/internal/config"
)

type Config = config.Config

type state struct {
	cfg			*Config
}

type command struct {
	name 		string
	arguments	[]string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Username is required")
	}

	if len(cmd.arguments) > 1 {
		return errors.New("Too many arguments, only username should be provided")
	}

	if err := s.cfg.SetUser(cmd.arguments[0]); err != nil {
		return err
	}
	
	fmt.Println("User has been set")

	return nil
}

func addConnectionString(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Connection string missing")
	}

	if len(cmd.arguments) > 1 {
		return errors.New("Too many arguments, only connection string should be provided")
	}

	if err := s.cfg.AddProtocol(cmd.arguments[0]); err != nil {
		return err
	}

	fmt.Println("Connection string has been set")

	return nil
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
	fmt.Println("BEFORE", c.list)
	if _, exists := c.list[name]; exists {
		return errors.New("command already exists")
	}
	
	c.list[name] = f

	fmt.Println("AFTER", c.list)

	return nil
}