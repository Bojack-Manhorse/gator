package cli

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/Bojack-Manhorse/pokedexcli/gator/internal/database"
	"github.com/Bojack-Manhorse/pokedexcli/gator/rss"
	"github.com/google/uuid"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	ctx := context.Background()
	return func(s *State, cmd Command) error {
		user, err := s.db.GetUser(ctx, s.config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
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

func Aggregate(s *State, cmd Command) error {

	/*
		feedUrl := "https://www.wagslane.dev/index.xml"

		feed, err := rss.FetchFeed(feedUrl)

		if err != nil {
			return err
		}

		fmt.Println(feed) */
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("Error, no time given.")
	}

	if len(cmd.Arguments) > 1 {
		return fmt.Errorf("Error, too many commands given.")
	}

	secBetweenReqs := cmd.Arguments[0]

	num, err := strconv.Atoi(secBetweenReqs)

	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Duration(num) * time.Second)

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Println("Error when scraping feed.")
		}
	}

	return nil
}

func scrapeFeeds(s *State) error {
	fmt.Println("Attempting to scrape feeds:")

	feed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return err
	}

	feedIdentifier := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: feed.ID,
	}

	err = s.db.MarkFeedFetched(context.Background(), feedIdentifier)

	if err != nil {
		return err
	}

	feedContent, err := rss.FetchFeed(feed.Url)

	if err != nil {
		return err
	}

	fmt.Println(feedContent.Channel.Title)
	fmt.Println(feedContent.Channel.Description)

	for _, item := range feedContent.Channel.Item {
		pubTimeLiteral, err := time.Parse(time.RFC1123Z, item.PubDate)

		var pubTimeSql sql.NullTime

		if err != nil {
			pubTimeSql = sql.NullTime{
				Time:  pubTimeLiteral,
				Valid: false,
			}
		} else {
			pubTimeSql = sql.NullTime{
				Time:  pubTimeLiteral,
				Valid: true,
			}
		}

		postStruct := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			Published:   pubTimeSql,
		}
		_, err = s.db.CreatePost(context.Background(), postStruct)

		if err != nil {
			fmt.Printf("Error adding post '%s' to database\n", item.Title)
		}
		fmt.Println(item)
	}

	return nil

}

func addFollowToUser(s *State, ctx context.Context, user string, feedURL string) (database.CreateFeedFollowRow, error) {

	userStruct, err := s.db.GetUser(ctx, user)

	if err != nil {
		return database.CreateFeedFollowRow{}, fmt.Errorf("Could not find current user in database")
	}

	feed, err := s.db.GetFeedFromURL(ctx, feedURL)

	if err != nil {
		return database.CreateFeedFollowRow{}, fmt.Errorf("Feed not found in database")
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userStruct.ID,
		FeedID:    feed.ID,
	}

	return s.db.CreateFeedFollow(ctx, feedFollow)
}

func AddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("No feed name given.")
	} else if len(cmd.Arguments) == 1 {
		return fmt.Errorf("No url given")
	} else if len(cmd.Arguments) > 2 {
		return fmt.Errorf("Too many arguments given")
	}

	feedName, url := cmd.Arguments[0], cmd.Arguments[1]

	ctx := context.Background()

	feedStruct := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       url,
		UserID:    user.ID,
	}

	result, err := s.db.CreateFeed(ctx, feedStruct)

	if err != nil {
		return fmt.Errorf("Could not insert entry into feeds table.")
	}

	fmt.Println(result)

	_, err = addFollowToUser(s, ctx, user.Name, url)

	if err != nil {
		return fmt.Errorf("Could not follow the feed '%s", feedName)
	}

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

func FollowFeed(s *State, cmd Command, user database.User) error {

	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("No url given")
	} else if len(cmd.Arguments) > 1 {
		return fmt.Errorf("Too many arguments given")
	}

	ctx := context.Background()

	feedURL := cmd.Arguments[0]

	data, err := addFollowToUser(s, ctx, user.Name, feedURL)

	if err != nil {
		return fmt.Errorf("Could not add follow to database")
	}

	fmt.Printf("User '%s' has followed the feed name '%s'\n", user.Name, data.FeedName)

	return nil
}

func GetAllFollowsOfUser(s *State, cmd Command, user database.User) error {

	if len(cmd.Arguments) > 0 {
		return fmt.Errorf("Too many arguments given")
	}

	ctx := context.Background()

	data, err := s.db.GetFeedFollowsForUser(ctx, user.Name)

	if err != nil {
		return fmt.Errorf("Culd not find user %s in database", user.Name)
	}

	fmt.Printf("User '%s' has followed the following feeds:\n", user.Name)
	for _, follow := range data {
		fmt.Println(follow)
	}

	return nil
}

func UnfollowFeed(s *State, cmd Command, user database.User) error {

	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("No feed given")
	} else if len(cmd.Arguments) > 1 {
		return fmt.Errorf("Too many arguments given")
	}

	feedUrl := cmd.Arguments[0]

	feedStruct, err := s.db.GetFeedFromURL(context.Background(), feedUrl)

	if err != nil {
		return err
	}

	followIdentifier := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedStruct.ID,
	}

	_, err = s.db.DeleteFeedFollow(context.Background(), followIdentifier)

	fmt.Printf("User %s has unfollowed feed %s \n", user.Name, feedStruct.Url)

	return nil
}

func Browse(s *State, cmd Command) error {
	var numberResults int

	if len(cmd.Arguments) > 1 {
		return fmt.Errorf("Too many arguments given")
	} else if len(cmd.Arguments) == 0 {
		fmt.Println("Defaulting to 2 results")
		numberResults = 2
	} else {
		val, err := strconv.Atoi(cmd.Arguments[0])

		if err != nil {
			fmt.Println("Defaulting to 2 results")
			numberResults = 2
		} else {
			numberResults = val
		}
	}

	results, err := s.db.GetPostsForUser(context.Background(), int32(numberResults))

	if err != nil {
		return err
	}

	for _, item := range results {
		fmt.Println(item)
	}

	return nil

}
