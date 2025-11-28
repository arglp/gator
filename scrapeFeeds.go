package main

import (
	"context"
	"errors"
	"time"
	"fmt"
	"database/sql"
	"github.com/arglp/gator/internal/database"
)

func scrapeFeeds(s* state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return errors.New("couldn't get next feed")
	}
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time: time.Now().UTC(), 
			Valid: true,
			},
		UpdatedAt: time.Now().UTC(),
		ID: feed.ID,
	})
	if err != nil {
		return errors.New("couldn't mark as fetched")
	}
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	fmt.Printf("channel title: %s\n", rssFeed.Channel.Title)
	fmt.Printf("channel link: %s\n", rssFeed.Channel.Link)
	fmt.Printf("channel description: %s\n",  rssFeed.Channel.Description)
	fmt.Println("items:")
	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("item title: %s\n", item.Title)
		fmt.Printf("item link: %s\n", item.Link)
		fmt.Printf("item description: %s\n", item.Description)
		fmt.Printf("item publication date: %s\n", item.PubDate)
	}
return nil
}
