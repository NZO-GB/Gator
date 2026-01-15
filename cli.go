package main

import (
	"fmt"
	"errors"
	"context"
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Username is required")
	}

	if len(cmd.arguments) > 1 {
		return errors.New("Too many arguments, only username should be provided")
	}

	username := cmd.arguments[0]

	if _, err := s.db.GetUser(context.Background(), username); err != nil {
		return err
	}

	if err := s.cfg.SetUser(username); err != nil {
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

	connectionString := cmd.arguments[0]

	if err := s.cfg.AddProtocol(connectionString); err != nil {
		return err
	}

	fmt.Println("Connection string has been set")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Username is required")
	}

	if len(cmd.arguments) > 1 {
		return errors.New("Too many arguments, only username should be provided")
	}

	username := cmd.arguments[0]
	userParams := s.createUserParams(username)

	user, err := s.db.CreateUser(context.Background(), userParams)

	if err != nil {
		return err
	}

	s.cfg.SetUser(username)
	
	fmt.Println("User has been created")
	fmt.Println(user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		return err
	}
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
	if _, exists := c.list[name]; exists {
		return errors.New("command already exists")
	}
	
	c.list[name] = f

	return nil
}