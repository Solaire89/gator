package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Solaire89/gator/internal/database"
	"github.com/google/uuid"
)

func printUser(user database.User) {
	fmt.Printf(" * ID:		%v\n", user.ID)
	fmt.Printf(" * Name:	%v\n", user.UserName)
}

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
	printUser(user)
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
	fmt.Println("Database reset successfully!")
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

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf(" * ID:			%v\n", feed.ID)
	fmt.Printf(" * Created at:	%v\n", feed.CreatedAt)
	fmt.Printf(" * Updated at:	%v\n", feed.UpdatedAt)
	fmt.Printf(" * Feed name:	%v\n", feed.FeedName)
	fmt.Printf(" * Feed url:	%v\n", feed.FeedUrl)
	fmt.Printf(" * User Name:	%v\n", user.UserName)
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
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
	feedFollowParams := database.CreateFeedFollowParams{
		ID:        fp.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	fmt.Println("Feed created successfully!")
	printFeed(feed, user)
	fmt.Println("")
	fmt.Println("Feed followed successfully!")
	printFeedFollow(feedFollow.FeedName, user.UserName)
	fmt.Println("=====================================")
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

func printFeedFollow(feed, user string) {
	fmt.Printf(" * Feed name:	%v\n", feed)
	fmt.Printf(" * User name:	%v\n", user)
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}
	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}
	createdFeedFollow, err := s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	printFeedFollow(createdFeedFollow.FeedName, createdFeedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}
	if len(feedFollows) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}
	for _, feed := range feedFollows {
		fmt.Printf("Feed name: %s\n", feed.FeedName)
	}
	return nil
}

func unfollowFeedCommand(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed url>", cmd.Name)
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	if err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return fmt.Errorf("couldn't delete feed follow: %w", err)
	}
	fmt.Printf("Successfully unfollowed %s\n", cmd.Args[0])
	return nil
}
