package cli

import (
	"fmt"

	"github.com/Bojack-Manhorse/pokedexcli/gator/internal/config"
)

type State struct {
	pointer *config.Config
}

type Command struct {
	name      string
	arguments []string
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("Error, command arguments is empty.")
	}

	if len(cmd.arguments) > 1 {
		return fmt.Errorf("Error, too many commands given.")
	}

	userName := cmd.arguments[0]

	return s.pointer.SetUser(userName)
}
