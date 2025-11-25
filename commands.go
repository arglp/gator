package main

import (
	"errors"
	"fmt"
	"context"
	"time"

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
		return errors.New("Unknown command")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) { 
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("argument is required")
	}

	name := cmd.args[0]
	
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return errors.New("user doesn't exist")
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	
	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("argument is required")
	}

	name := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return errors.New("username already exists")
	}

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: cmd.args[0],
	})

	if err != nil {
		return errors.New("couldn't register user")
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}

	fmt.Printf("new user %s was created", name)
	fmt.Println(user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return errors.New("couldn't delete")
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return errors.New ("couldn't get users")
	}
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
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return errors.New("couldn't fetch feed")
	}
	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return errors.New("required more arguments")
	}

	user, err := s.db.GetUser(context.Background(),s.cfg.CurrentUserName)
	if err != nil {
		return errors.New("couldn't get current user")
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
		UserID : userId,
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

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return errors.New("required more arguments")
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return errors.New ("couldn't find feed")
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return errors.New ("couldn't find user")
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

func handlerFollowing( s* state, cmd command) error {
	user, err := s.db.GetUser(context.Background(),s.cfg.CurrentUserName)
	if err != nil {
		return errors.New("couldn't find user")
	}
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
