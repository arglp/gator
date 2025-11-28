package main

import (
	"github.com/arglp/gator/internal/database"
	"context"
	"errors"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return errors.New ("couldn't find user")
		}
		err = handler(s, cmd, user)
		if err != nil {
			return err
		}
		return nil
	}
}