package commandss

import (
	"context"
	"database/sql"
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
	if len(cmd.Args) < 1 {
		return fmt.Errorf("timout between fetches is required, eg: 5s, 2m, 1h")
	}
	timeout, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration format: %w", err)
	}
	ticker := time.NewTicker(timeout)
	for ; ; <-ticker.C {
		fmt.Printf("Collecting feeds every %s\n", timeout)
		scrapeFeed(s)
	}
}

func HandlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("feed URL are required")
	}
	ctx := context.Background()

	feed, err := s.Db.GetFeedsUrl(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error retrieving feed by URL: %w", err)
	}
	if _, err := s.Db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}); err != nil {
		return fmt.Errorf("error adding feed follow: %w", err)
	}
	fmt.Printf("Successfully followed feed: %s\n", feed.Name)

	return nil

}

func HandlerFollowing(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("no arguments expected for following command")
	}
	ctx := context.Background()

	following, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error retrieving followed feeds: %w", err)
	}
	if len(following) == 0 {
		fmt.Println("You are not following any feeds.")
	}
	for _, follow := range following {
		fmt.Printf("%s\n", follow.FeedName)
	}
	return nil

}

func HandlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) <= 1 {
		return fmt.Errorf("feed URL and name are required")
	}

	if feed, err := s.Db.AddFeed(context.Background(), database.AddFeedParams{
		ID:        uuid.New(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}); err != nil {
		return fmt.Errorf("error adding feed: %w", err)
	} else {
		fmt.Printf("%+v\n", feed)
	}
	if err := HandlerFollow(s, Command{
		Name: "follow",
		Args: []string{cmd.Args[1]},
	}, user); err != nil {
		return fmt.Errorf("error following feed after adding: %w", err)
	}
	fmt.Printf("Feed %s added and followed successfully.\n", cmd.Args[0])

	return nil
}

func HandlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("feed URL")
	}
	data, err := s.Db.GetFeedsUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error retrieving feed by URL: %w", err)
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("only one feed ID is expected")
	}
	if err := s.Db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: data.ID,
	}); err != nil {
		return fmt.Errorf("error unfollowing feed: %w", err)
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

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(s *State, cmd Command) error {
	return func(s *State, cmd Command) error {
		data, err := s.Db.GetUser(context.Background(), s.ConfigFile.CURRENT_USER)
		if err != nil {
			return fmt.Errorf("error retrieving current user: %w", err)
		}
		result := handler(s, cmd, data)
		if result != nil {
			return fmt.Errorf("error in handler: %w", result)
		}
		return nil

	}

}

func scrapeFeed(s *State) error {
	ctx := context.Background()
	feed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("error retrieving next feed to fetch: %w", err)
	}
	for _, f := range feed {
		s.Db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
			LastFetechAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: time.Now(),
			ID:        f.ID,
		})
		rss_feed, err := rss.FetchFeed(ctx, f.Url)
		if err != nil {
			return fmt.Errorf("error fetching feed: %w", err)
		}
		for _, item := range rss_feed.Channel.Items {
			s.Db.CreatePost(ctx, database.CreatePostParams{
				ID:        uuid.New(),
				UpdatedAt: time.Now(),
				Title:     item.Title,
				Url:       item.Link,
				FeedID:    f.ID,
				Description: sql.NullString{
					String: item.Description,
					Valid:  true,
				},
				PublishedAt: PubDate(item.PubDate),
			})
		}
	}
	return nil
}
