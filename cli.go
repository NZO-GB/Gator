package main

import (
	"errors"
	"time"
	uuid "github.com/google/uuid"
	database "github.com/NZO-GB/Gator/internal/database"
	config "github.com/NZO-GB/Gator/internal/config"
)

type Config = config.Config

type state struct {
	cfg			*Config
	db			*database.Queries
}

func (s state) createUserParams(name string) database.CreateUserParams {
	userID := uuid.New()
	now := time.Now()

	user := database.CreateUserParams{
		ID:			userID,
		CreatedAt:  now,
		UpdatedAt:  now,
		Name:		name,
	}

	return user
}

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