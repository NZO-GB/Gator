package main

import (
	"time"
	"context"
	"fmt"
	uuid "github.com/google/uuid"
	database "github.com/NZO-GB/Gator/internal/database"
	config "github.com/NZO-GB/Gator/internal/config"
)

type state struct {
	cfg			*config.Config
	db			*database.Queries
}

func (s state) CreateUserParams(name string) database.CreateUserParams {
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

func (s state) CreateFeedParams(user database.User, feedname string, url string) (database.CreateFeedParams, error) {

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

func (s state) CreateCreatePostParams(
	title string, url string, description string, publishedAt time.Time, feedID uuid.UUID) (database.CreatePostParams, error) {

	postParams := database.CreatePostParams{
	ID:          uuid.New(),
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
	Title:       title,
	Url:         url,
	Description: description,
	PublishedAt: publishedAt,
	FeedID:      feedID,
	}

	return postParams, nil

}

func scrapeFeeds(s *state) error {

	fmt.Printf("---------\nScraping feeds...\n---------\n")

	feedQuery, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	url := feedQuery.Url

	err = s.db.MarkFeedFetched(context.Background(), url)
	if err != nil {
		return err
	}

	RSSfeed, err := fetchFeed(context.Background(), url)
	if err != nil {
		fmt.Println("Error fetching feed")
		return err
	}

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	for _, item := range(RSSfeed.Channel.Item) {
		err := createPostHelper(s, item, feed.ID)
		if err != nil {
			return err
		}
	}

	fmt.Printf("---------\nFeeds Scraped Succesfully\n---------\n")
	
	return nil
}