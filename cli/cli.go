package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/Bojack-Manhorse/pokedexcli/gator/internal/config"
	"github.com/Bojack-Manhorse/pokedexcli/gator/internal/database"
	"github.com/google/uuid"
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

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("Error, command arguments is empty.")
	}

	if len(cmd.Arguments) > 1 {
		return fmt.Errorf("Error, too many commands given.")
	}

	userName := cmd.Arguments[0]

	ctx := context.Background()

	_, err := s.db.GetUser(ctx, userName)

	if err != nil {
		return fmt.Errorf("User %s not found in database.", userName)
	}

	return s.config.SetUser(userName)
}

func Register(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("Error, command arguments is empty.")
	}

	if len(cmd.Arguments) > 1 {
		return fmt.Errorf("Error, too many commands given.")
	}

	ctx := context.Background()

	userStruct := database.CreateUserParams{
		Name:      cmd.Arguments[0],
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.CreateUser(ctx, userStruct)

	if err != nil {
		return err
	}

	err = HandlerLogin(s, cmd)

	return err

}

func Reset(s *State, cmd Command) error {
	ctx := context.Background()
	return s.db.Reset(ctx)
}

func ListUsers(s *State, cmd Command) error {
	ctx := context.Background()
	listOfUsers, err := s.db.GetUsers(ctx)

	if err != nil {
		return err
	}

	currentUser := s.config.CurrentUserName

	for _, name := range listOfUsers {
		if name == currentUser {
			fmt.Printf("* %s (current)\n", name)
		} else {
			fmt.Printf("* %s\n", name)
		}
	}

	return nil
}
