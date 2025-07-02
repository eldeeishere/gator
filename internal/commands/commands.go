package commandss

import (
	"context"
	"fmt"
	"time"

	"github.com/eldeeishere/gator/internal/config"
	"github.com/eldeeishere/gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	Db         *database.Queries
	ConfigFile config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required")
	}
	ctx := context.Background()
	if user, err := s.Db.GetUser(ctx, cmd.Args[0]); err != nil {
		if cmd.Args[0] != user.Name {
			return fmt.Errorf("error creating user: %w", err)
		}
		return fmt.Errorf("error retrieving user: %w", err)
	} else {
		fmt.Printf("Welcome back, %s!\n", user.Name)
	}
	if err := s.ConfigFile.SetUser(cmd.Args[0]); err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("User %s logged in successfully.\n", cmd.Args[0])
	return nil
}

func HandlerUsers(s *State, cmd Command) error {
	ctx := context.Background()
	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving users: %w", err)
	}
	if len(users) == 0 {
		fmt.Println("No users found.")
		return nil
	}
	for _, user := range users {
		if user == s.ConfigFile.CURRENT_USER {
			fmt.Printf("* %s (current)\n", user)
			continue
		}
		fmt.Printf("* %s\n", user)
	}
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("username is required")
	}
	ctx := context.Background()
	arg := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	if user, err := s.Db.CreateUser(ctx, arg); err != nil {
		return fmt.Errorf("Error creating user: %w\n", user.Name)
	} else {
		fmt.Printf("User %s registered successfully.\n", user.Name)
		if err := s.ConfigFile.SetUser(user.Name); err != nil {
			return fmt.Errorf("failed to set user: %w", err)
		}
	}
	fmt.Printf("Debug: User %s created with ID %s\n", arg.Name, arg.ID)
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	ctx := context.Background()
	rss.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")

	return nil
}

func HandlerReset(s *State, cmd Command) error {
	ctx := context.Background()
	if err := s.Db.Reset(ctx); err != nil {
		return fmt.Errorf("error resetting database: %w", err)
	}
	fmt.Println("Database reset successfully.")
	return nil

}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, exists := c.Commands[cmd.Name]; exists {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}

func (c *Commands) Register(name string, handler func(*State, Command) error) {
	if c.Commands == nil {
		c.Commands = make(map[string]func(*State, Command) error)
	}
	c.Commands[name] = handler
}
