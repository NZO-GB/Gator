package main

import(
	"context"
	"time"
	"fmt"
	database "github.com/NZO-GB/Gator/internal/database"
	uuid "github.com/google/uuid"
)

func middlewareLoggedIn(
	handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
		return func(s *state, cmd command) error {
			username := s.cfg.CurrentUserName
			user, err := s.db.GetUser(context.Background(), username)
			if err != nil {
				return err
			}
			return handler(s, cmd, user)
		}
	}


func createPostHelper(s *state, item RSSItem, feedID uuid.UUID) error {

	const layout = time.RFC1123Z

	pubDate, err := time.Parse(layout, item.PubDate)
	if err != nil {
		fmt.Println("Error parsing published time")
		return err
	}

	postParams, err := s.CreateCreatePostParams(item.Title, item.Link, item.Description, pubDate, feedID)
	s.db.CreatePost(context.Background(), postParams)

	fmt.Println(item.Title)

	return nil
}