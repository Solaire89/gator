package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Solaire89/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s <name>\n", cmd.Name)
	}
	name := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("error logging in: %w", err)
	}
	fmt.Printf("username set to %v\n", cmd.Args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>\n", cmd.Name)
	}

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
	}
	user, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Printf("successfully set user to %s\n", user.Name)
	log.Printf("user set to %+v", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s <name>\n", cmd.Name)
	}
	if err := s.db.ClearTable(context.Background()); err != nil {
		return fmt.Errorf("error: %w", err)
	}
	return nil
}
