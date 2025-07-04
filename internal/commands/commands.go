package commandss

import (
	"context"
	"fmt"
	"time"

	"github.com/eldeeishere/gator/internal/config"
	"github.com/eldeeishere/gator/internal/database"
	"github.com/eldeeishere/gator/internal/rss/rss"
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
	user, err := s.Db.GetUser(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("user %s does not exist: %w", cmd.Args[0], err)
	}
	fmt.Printf("Welcome back, %s!\n", user.Name)

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
	user, err := s.Db.CreateUser(ctx, arg)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err) // Zmeň z user.Name na err
	}
	fmt.Printf("User %s registered successfully.\n", user.Name)
	if err := s.ConfigFile.SetUser(user.Name); err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("Debug: User %s created with ID %s\n", arg.Name, arg.ID)
	return nil
}

func HandlerAgg(s *State, cmd Command) error {
	ctx := context.Background()
	content, err := rss.FetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("error fetching RSS feed: %w", err)
	}
	if len(content.Channel.Items) == 0 {
		fmt.Println("No items found in the RSS feed.")
	}
	fmt.Printf("Feed: %+v\n", content)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) <= 1 {
		return fmt.Errorf("feed URL and name are required")
	}
	userID, err := s.getCurrentUserUUID()
	if err != nil {
		return fmt.Errorf("error retrieving current user UUID: %w", err)
	}
	if feed, err := s.Db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    userID,
	}); err != nil {
		return fmt.Errorf("error adding feed: %w", err)
	} else {
		fmt.Printf("%+v\n", feed)
	}

	return nil
}

func HandlerFeeds(s *State, cmd Command) error {

	feed, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error retrieving feeds: %w", err)
	}
	if len(feed) == 0 {
		fmt.Println("No feeds found.")
	}
	for _, f := range feed {
		fmt.Printf("%s\n%s\n", f.Name, f.Name_2)
	}

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

func (s *State) getCurrentUserUUID() (uuid.UUID, error) {
	if s.ConfigFile.CURRENT_USER == "" {
		return uuid.Nil, fmt.Errorf("no current user set in config")
	}
	data, err := s.Db.GetUser(context.Background(), s.ConfigFile.CURRENT_USER)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error retrieving current user: %w", err)
	}
	return data.ID, nil

}
