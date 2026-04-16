package main

import (
	"errors"
	"time"
	"context"
	"net/http"
	"io"
	"encoding/xml"
	uuid "github.com/google/uuid"
	database "github.com/NZO-GB/Gator/internal/database"
	config "github.com/NZO-GB/Gator/internal/config"
)

var client = http.Client{}

type state struct {
	cfg			*config.Config
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

func (s state) createFeedParams(user database.User, feedname string, url string) (database.CreateFeedParams, error) {

	feedParams := database.CreateFeedParams {
		ID:			uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Name:		feedname,
		Url: 		url,
		UserID:		user.ID,
	}

	return feedParams, nil
}

func (s state) CreateFeedFollowParams(user database.User, url string) (database.CreateFeedFollowParams, error) {

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return database.CreateFeedFollowParams{}, err
	}

	feedFollow := database.CreateFeedFollowParams {
		ID:			uuid.New(),
		UserID:		user.ID,
		FeedID: 	feed.ID,
		CreatedAt: 	time.Now(),
		UpdatedAt: 	time.Now(),
	}

	return feedFollow, nil
}

func (s state) CreateRemoveFeedFollowParams(user database.User, url string) (database.RemoveFeedFollowParams, error) {

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return database.RemoveFeedFollowParams{}, err
	}

	removeFollow := database.RemoveFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	return removeFollow, nil
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

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)

	var xmlStruct RSSFeed

	if err := xml.Unmarshal(data, &xmlStruct); err != nil {
		return nil, err
	}

	return &xmlStruct, err
}
