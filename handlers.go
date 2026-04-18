package main

import (
	"context"
	"errors"
	"fmt"
	database "github.com/NZO-GB/Gator/internal/database"
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
	
	fmt.Printf("User has been set: %s\n", username)

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

	if err = s.cfg.SetUser(username); err != nil {
		return err
	}
	
	fmt.Printf("User has been created: %s \n", user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.Reset(context.Background()); err != nil {
		return err
	}
	fmt.Println("Database reset succesful")
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
		fmt.Printf("Got user: %s \n", printName)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {

	url := cmd.arguments[0]

	_, err := fetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	s.scrapeFeeds()
	
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return errors.New("Name and URL are required")
	}

	feedname := cmd.arguments[0]
	url := cmd.arguments[1]

	feedParams, err := s.createFeedParams(user, feedname, url)
	if err != nil {
		fmt.Println("Error creating params")
		return err
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		fmt.Println("Error creating feed")
		return err
	}

	fmt.Printf("Added feed: '%s' \n", feed.Name)

	cmdFeedFollow := command {
		arguments: []string{feed.Url},
	}

	if err := handlerFollow(s, cmdFeedFollow, user); err != nil {
		fmt.Println("Error following")
		return err
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, f := range feeds { // XXXXX Do we need to print User???? Can we even do that?
		userID := f.UserID
		username, err := s.db.GetUserByID(context.Background(), userID)
		if err != nil {
			return err
		}
		fmt.Printf("Feed: %s \n", f.Name)
		fmt.Printf("URL: %s \n ", f.Url)
		fmt.Printf("User: %s \n", username)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {

	url := cmd.arguments[0]

	feedFollowParams, err := s.CreateFeedFollowParams(user, url)
	if err != nil {
		fmt.Println("Error creating params")
		return err
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		fmt.Println("Error creating follow")
		return err
	}

	fmt.Printf("Added feed: '%s' to User: '%s' \n", feedFollow.FeedName, feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	feedfollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("---\nPrinting feeds for '%s': \n", user.Name)

	for _, feed := range(feedfollows) {
		fmt.Printf("Feed: %s\n", feed.FeedName)
	}

	fmt.Println("---")

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {

	url := cmd.arguments[0]

	removeFeedFollowParams, err := s.CreateRemoveFeedFollowParams(user, url)
	if err != nil {
		fmt.Println("Error creating params")
		return err
	}

	err = s.db.RemoveFeedFollow(context.Background(), removeFeedFollowParams)
	if err != nil {
		fmt.Println("Error unfollowing")
		return err
	}

	fmt.Printf("Feed: '%s' has been deleted for User: '%s'\n", url, user.Name)

	return nil

}
