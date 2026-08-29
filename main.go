package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Bojack-Manhorse/pokedexcli/gator/cli"
	"github.com/Bojack-Manhorse/pokedexcli/gator/internal/config"
	"github.com/Bojack-Manhorse/pokedexcli/gator/internal/database"

	_ "github.com/lib/pq"
)

const dbUrl string = "DATABASE_URL/gator?sslmode=disable"

type CommandHandlers map[string]func(*cli.State, cli.Command) error

var commandMapping = CommandHandlers{
	"login":    cli.HandlerLogin,
	"register": cli.Register,
	"reset":    cli.Reset,
	"users":    cli.ListUsers,
}

func getCurrentConfig() config.Config {
	currentConfig, err := config.Read()

	if err != nil {
		fmt.Println("Cannot read config file due to error:", err)
		os.Exit(1)
	}

	return currentConfig
}

func connectToDatabase(url string) *database.Queries {
	db, err := sql.Open("postgres", url)

	if err != nil {
		fmt.Printf("Error when opening database with url %s: %v.\n", url, err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	return dbQueries
}

func initState() cli.State {
	var storedState cli.State

	currentConfig := getCurrentConfig()
	dbQueries := connectToDatabase(dbUrl)

	storedState.SetConfig(&currentConfig)
	storedState.SetDb(dbQueries)

	return storedState
}

func initCommands(commandMap CommandHandlers) cli.Commands {
	var commands cli.Commands

	commands.CommandMap = make(CommandHandlers)

	for key, value := range commandMapping {
		err := (&commands).Register(key, value)
		if err != nil {
			fmt.Printf("Error when registering command %s, %v\n", key, err)
			os.Exit(1)
		}
	}

	return commands
}

func parseArgs() (string, []string) {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Too few arguments given.")
		os.Exit(1)
	}

	commandName := args[1]

	commandArguments := args[2:]

	return commandName, commandArguments
}

func runCommand(commandName string, commandArgs []string, state cli.State, commandList cli.Commands) {

	var command cli.Command = cli.Command{
		Name:      commandName,
		Arguments: commandArgs,
	}

	err := commandList.Run(&state, command)

	if err != nil {
		fmt.Printf("Error when running command %v\n", err)
		os.Exit(1)
	}
}

func main() {

	storedState := initState()
	eligibleCommands := initCommands(commandMapping)
	commandName, commandArguments := parseArgs()

	runCommand(commandName, commandArguments, storedState, eligibleCommands)

}
