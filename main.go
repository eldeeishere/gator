package main

import (
	"fmt"
	"os"

	commandss "github.com/eldeeishere/gator/internal/commands"
	"github.com/eldeeishere/gator/internal/config"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No command provided. Please provide a command to execute.")
		os.Exit(1)
	}
	configFile, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}
	state := commandss.State{
		ConfigFile: configFile,
	}

	commands := commandss.Commands{
		Commands: make(map[string]func(*commandss.State, commandss.Command) error),
	}
	commands.Register("login", commandss.HandlerLogin)

	cmd := commandss.Command{
		Name: args[1],
		Args: args[2:],
	}
	if err := commands.Run(&state, cmd); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
