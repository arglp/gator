package main

import (
	"errors"
	"fmt"
	"context"
	"time"
	"strconv"

	"github.com/google/uuid"
	"github.com/arglp/gator/internal/database"
)

type command struct {
	name 	string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New("unknown command")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) { 
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("please provide a user name")
	}

	name := cmd.args[0]
	
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user %s doesn't exist", name)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	
	fmt.Printf("%s has been set as active user", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("please provide a username to register")
	}

	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("username %s already exists", name)
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: cmd.args[0],
	})

	if err != nil {
		return fmt.Errorf("couldn't register user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set user: %w",err)
	}

	fmt.Printf("registered new user %s\n", user.Name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("no users registered")
	}

	fmt.Println("registered users:")

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Printf("%s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("please enter a time string")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("collecting feeds every %.0fm%.0fs", timeBetweenRequests.Minutes(), timeBetweenRequests.Seconds())

	ticker := time.NewTicker(timeBetweenRequests)
	for ;; <-ticker.C {
		err = scrapeFeeds(s)
		fmt.Println("collecting feeds")
		if err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return errors.New("please provide name and url")
	}

	userId := user.ID
	name := cmd.args[0]
	url := cmd.args[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: name,
		Url: url,
		UserID: userId,
	})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Println("added feed:")
	fmt.Println(feed.ID)
	fmt.Println(feed.CreatedAt)
	fmt.Println(feed.UpdatedAt)
	fmt.Println(feed.Name)
	fmt.Println(feed.Url)
	fmt.Println(feed.UserID)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return errors.New("couldn't get feeds")
	}
	for _, feed := range feeds{
		fmt.Printf("name: %s, url: %s, user: %s\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("required more arguments")
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return errors.New ("couldn't find feed")
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return errors.New("couldn't follow")
	}

	fmt.Println("following new feed")
	fmt.Printf("user: %s, feed: %s\n", follow.UserName, follow.FeedName)
	return nil
}

func handlerFollowing( s* state, cmd command, user database.User) error {
	
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return errors.New("couldn't find followed feeds")
	}
	fmt.Printf("user: %s is following these feeds:\n", user.Name)
	for _, follow := range follows {
		fmt.Println(follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s* state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return errors.New("required more arguments")
	}
	url := cmd.args[0]

	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return err
	}
	err = s.db.DeleteFeedFollow(context.Background(),database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func handlerBrowse(s* state, cmd command, user database.User) error {
	
	var limit int32 = 2
	if len(cmd.args) > 0 {
		limit64, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err == nil {
			limit = int32(limit64)
		}
	}

	posts, err := s.db.GetPostForUser(context.Background(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit: limit,
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("title: %v\n", post.Title)
		fmt.Printf("link: %s\n", post.Url)
		fmt.Printf("item description: %v\n", post.Description)
		fmt.Printf("item publication date: %v\n", post.PublishedAt)
	}
return nil
}
