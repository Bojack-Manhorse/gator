package cli

import (
	"fmt"

	"github.com/Bojack-Manhorse/pokedexcli/gator/internal/config"
	"github.com/Bojack-Manhorse/pokedexcli/gator/internal/database"
)

type State struct {
	db     *database.Queries
	config *config.Config
}

type Command struct {
	Name      string
	Arguments []string
}

func (s *State) SetConfig(configPointer *config.Config) {
	s.config = configPointer
}

func (s *State) SetDb(dbPointer *database.Queries) {
	s.db = dbPointer
}

type Commands struct {
	CommandMap map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	v, exists := c.CommandMap[cmd.Name]

	if !exists {
		return fmt.Errorf("Error, command %s does not exist", cmd.Name)
	}

	return v(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) error {
	_, exist := c.CommandMap[name]

	if exist {
		return fmt.Errorf("Error, command %s already exists.", name)
	}

	c.CommandMap[name] = f

	return nil
}
