package main

import (
	"github.com/arglp/gator/internal/config"
	"errors"
	"fmt"
)

type state struct {
	cfg 	*config.Config
}

type command struct {
	name	string
	arguments	[]string
}

type commands struct {
	library map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.library[cmd.name]
	if !ok {
		return errors.New("Command not found")
	}
	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.library[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("No argument")
	}
	err := s.cfg.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}