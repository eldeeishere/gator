package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	commandss "github.com/eldeeishere/gator/internal/commands"
	"github.com/eldeeishere/gator/internal/config"
	"github.com/eldeeishere/gator/internal/database"
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
	dbUrl := configFile.DB_URL
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	dbQueries := database.New(db)
	state := commandss.State{
		Db:         dbQueries,
		ConfigFile: configFile,
	}

	commands := commandss.Commands{
		Commands: make(map[string]func(*commandss.State, commandss.Command) error),
	}
	commands.Register("login", commandss.HandlerLogin)
	commands.Register("register", commandss.HandlerRegister)
	commands.Register("reset", commandss.HandlerReset)
	commands.Register("users", commandss.HandlerUsers)
	commands.Register("agg", commandss.HandlerUsers)

	cmd := commandss.Command{
		Name: args[1],
		Args: args[2:],
	}
	if err := commands.Run(&state, cmd); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
