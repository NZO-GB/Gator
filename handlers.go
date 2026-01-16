package main

import(
	"errors"
	"fmt"
	"context"
)

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
	
	fmt.Println("User has been created:")
	fmt.Println(user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		return err
	}
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, u := range(users) {

		printName := "* " + u.Name
		if u.Name == s.cfg.CurrentUserName {
			printName += " (current)"
		}
		fmt.Println(printName)
	}

	return nil
}