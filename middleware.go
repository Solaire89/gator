package main

import (
	"context"
	"fmt"

	"github.com/Solaire89/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	handlerWrapper := func(s *state, cmd command) error {
		currentUser, err := s.db.GetUser(context.Background(), s.cfg.UserName)
		if err != nil {
			return fmt.Errorf("error getting user: %w", err)
		}
		return handler(s, cmd, currentUser)
	}
	return handlerWrapper
}
