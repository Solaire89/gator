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
	if len(cmd.Args) != 1 {
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
		UserName:  cmd.Args[0],
	}
	user, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = s.cfg.SetUser(user.UserName)
	if err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	fmt.Printf("successfully set user to %s\n", user.UserName)
	log.Printf("user set to %+v", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	if err := s.db.ClearTable(context.Background()); err != nil {
		return fmt.Errorf("error: %w", err)
	}
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	for _, user := range users {
		if user.UserName == s.cfg.UserName {
			fmt.Printf("* %s (current)\n", user.UserName)
		} else {
			fmt.Printf("* %s\n", user.UserName)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}
	user, err := s.db.GetUser(context.Background(), s.cfg.UserName)
	if err != nil {
		return fmt.Errorf("couldn't get user: %w", err)
	}
	fp := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		FeedName:  cmd.Args[0],
		FeedUrl:   cmd.Args[1],
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeed(context.Background(), fp)
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}
	fmt.Printf("Added feed: %+v", feed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	users, err := s.db.UserWhoCreatedFeed(context.Background())
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	for _, user := range users {
		fmt.Printf("Feed Name: %s\n", user.FeedName)
		fmt.Printf("Feed URL: %s\n", user.FeedUrl)
		fmt.Printf("User who created feed: %s\n", user.UserName)
	}

	return nil
}
