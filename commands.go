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