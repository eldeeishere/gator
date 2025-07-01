package commandss

import (
	"fmt"

	"github.com/eldeeishere/gator/internal/config"
)

type State struct {
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
	if err := s.ConfigFile.SetUser(cmd.Args[0]); err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Printf("User %s logged in successfully.\n", cmd.Args[0])
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
