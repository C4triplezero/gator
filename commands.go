package main

import (
	"context"
	"fmt"
	"time"

	"github.com/C4triplezero/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("couldn't created user: %w", err)
	}

	err = s.Config.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	_, err := s.DB.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user: %s does not exist", name)
	}

	err = s.Config.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User: %s has been set\n", name)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.DB.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't reset table: %w", err)
	}

	fmt.Println("Users table has been reset")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("unable to retrieve users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}

		fmt.Println("* " + user.Name)
	}

	return nil
}

func handlerAggregation(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("unable to fetch feed: %w", err)
	}

	fmt.Println(feed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("requires user to be logged in: %w", err)
	}

	feed, err := s.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("unable to create feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
	fmt.Printf("* User:          %s\n", user.Name)
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("unable to retrieve feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.DB.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("=====================================")
	}

	return nil
}
