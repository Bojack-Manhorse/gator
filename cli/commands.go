package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/Bojack-Manhorse/pokedexcli/gator/internal/database"
	"github.com/Bojack-Manhorse/pokedexcli/gator/rss"
	"github.com/google/uuid"
)

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

func Aggregate(s *State, cmd Command) error {
	feedUrl := "https://www.wagslane.dev/index.xml"

	feed, err := rss.FetchFeed(feedUrl)

	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func AddFeed(s *State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("No feed name given.")
	}

	if len(cmd.Arguments) == 1 {
		return fmt.Errorf("No url given")
	}

	if len(cmd.Arguments) > 2 {
		return fmt.Errorf("Too many arguments given")
	}

	feedName, url := cmd.Arguments[0], cmd.Arguments[1]

	ctx := context.Background()

	currentUser := s.config.CurrentUserName

	currentUserEntry, err := s.db.GetUser(ctx, currentUser)

	if err != nil {
		return fmt.Errorf("User %s does not exist in database.", currentUser)
	}

	feedStruct := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       url,
		UserID:    currentUserEntry.ID,
	}

	result, err := s.db.CreateFeed(ctx, feedStruct)

	if err != nil {
		return fmt.Errorf("Could not insert entry into feeds table.")
	}

	fmt.Println(result)

	return nil
}

func ListFeeds(s *State, cmd Command) error {
	ctx := context.Background()
	feeds, err := s.db.ListAllFeeds(ctx)

	if err != nil {
		return fmt.Errorf("Could not retrieve feeds from table.")
	}

	for _, feed := range feeds {
		fmt.Println(feed)
	}

	return nil
}
