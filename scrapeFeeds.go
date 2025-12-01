package main

import (
	"context"
	"errors"
	"time"
	"database/sql"

	"github.com/google/uuid"
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

	for _, item := range rssFeed.Channel.Item {

		title := sql.NullString{}
		if item.Title == "" {
			title.Valid = false
		} else {
			title.String = item.Title
			title.Valid = true
		}
		description := sql.NullString{}
		if item.Description == "" {
			description.Valid = false
		} else {
			description.String = item.Description
			description.Valid = true
		}

		timeLayouts := []string{
						time.RFC3339, 
						time.RFC822, 
						time.RFC850,
						time.RFC1123,
						time.ANSIC,
						time.UnixDate,
						time.RubyDate}

		pubDateParsed := time.Now()
		for _, layout := range timeLayouts {
			pubDateParsed, err = time.Parse(layout, item.PubDate)
			if err == nil {
				break
			}
		}

		publishedAt := sql.NullTime{}
		if item.PubDate == "" {
			publishedAt.Valid = false
		} else {
			publishedAt.Time = pubDateParsed
			publishedAt.Valid = true
		}

		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:			uuid.New(),
			CreatedAt: 	time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
			Title:      title,
			Url:		item.Link,
			Description: description,
			PublishedAt:	publishedAt,
			FeedID:		feed.ID,
		})
		if err != nil {
			return err
		}
	}
return nil
}
