package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/halfdan87/boot-go-blog-aggregator-2.0/internal/config"
	"github.com/halfdan87/boot-go-blog-aggregator-2.0/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Starting...")

	cfg, err := config.Read()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbQueries := database.New(db)

	s := &state{cfg: cfg, db: dbQueries}
	c := &commands{mapping: make(map[string]func(*state, command) error)}
	c.register("login", login)
	c.register("register", register)
	c.register("reset", deleteAll)
	c.register("users", list)
	c.register("agg", aggregate)
	c.register("addfeed", middlewareLoggedIn(addFeed))
	c.register("feeds", allFeeds)
	c.register("follow", middlewareLoggedIn(follow))
	c.register("following", middlewareLoggedIn(following))
	c.register("unfollow", middlewareLoggedIn(unfollow))

	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("No command given")
		os.Exit(1)
	}

	err = c.run(s, command{name: args[0], args: args[1:]})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	/*
		for {
			fmt.Print("Command: ")
			cmd, err := readCommand()
			if err != nil {
				panic(err)
			}
			err = c.run(s, cmd)
			if err != nil {
				fmt.Println(err)
			}
		}
	*/
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserId)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}

func list(s *state, _ command) error {
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserId {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func deleteAll(s *state, _ command) error {
	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	return s.cfg.SetUser("")
}

func login(s *state, cmd command) error {
	// get user from db and set it in config
	if len(cmd.args) == 0 {
		fmt.Println("No username given")
		os.Exit(1)
	}
	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return err
	}
	return s.cfg.SetUser(username)
}

func aggregate(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Println("No url given")
		//os.Exit(1)
	}
	//url := cmd.args[0]

	url := "https://www.wagslane.dev/index.xml"
	feed, err := fetchRssFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Printf("Feed: %v\n", feed)
	return nil
}

func allFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		// find user who created this feed
		dbUser, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}

		fmt.Printf("* %s, %s, %s\n", feed.Name, feed.Url, dbUser.Name)
	}
	return nil
}

func follow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		fmt.Println("No url given")
		os.Exit(1)
	}
	url := cmd.args[0]

	dbFeed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	createParams := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		FeedID: dbFeed.ID,
		UserID: user.ID,
	}

	dbFeedFollow, err := s.db.CreateFeedFollow(context.Background(), createParams)
	if err != nil {
		return err
	}
	fmt.Printf("Feed follow created: %s - %s\n", dbFeedFollow.FeedName, dbFeedFollow.UserName)
	return nil
}

func unfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		fmt.Println("No url given")
		os.Exit(1)
	}
	url := cmd.args[0]

	dbFeed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	deleteParams := database.DeleteFeedFollowByUserIdAndFeedUrlParams{
		Url:    url,
		UserID: user.ID,
	}

	err = s.db.DeleteFeedFollowByUserIdAndFeedUrl(context.Background(), deleteParams)
	if err != nil {
		return err
	}
	fmt.Printf("Feed unfollowed: %s - %s\n", dbFeed.Name, dbFeed.Url)
	return nil
}

func following(s *state, _ command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}
	return nil
}

func addFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		fmt.Println("No url given")
		os.Exit(1)
	}
	name := cmd.args[0]
	url := cmd.args[1]

	createParams := database.CreateFeedParams{
		ID:     uuid.New(),
		Name:   name,
		Url:    url,
		UserID: user.ID,
	}

	dbFeed, err := s.db.CreateFeed(context.Background(), createParams)
	if err != nil {
		return err
	}

	fmt.Printf("Feed created %v\n", dbFeed)

	createFollowParams := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		FeedID: dbFeed.ID,
		UserID: user.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), createFollowParams)
	if err != nil {
		return err
	}

	return nil
}

func register(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		fmt.Println("No username given")
		os.Exit(1)
	}
	name := cmd.args[0]
	createParams := database.CreateUserParams{
		Name: name,
		ID:   uuid.New(),
	}

	dbUser, err := s.db.CreateUser(context.Background(), createParams)
	if err != nil {
		return err
	}
	fmt.Printf("User created %v\n", dbUser)
	return s.cfg.SetUser(name)
}

func readCommand() (command, error) {
	var cmd command
	var err error
	cmd.name, err = readString()
	if err != nil {
		return cmd, err
	}
	cmd.args, err = readArgs()
	return cmd, err
}

func readString() (string, error) {
	var s string
	_, err := fmt.Scanln(&s)
	return s, err
}

func readArgs() ([]string, error) {
	var args []string
	for {
		var arg string
		_, err := fmt.Scanln(&arg)
		if err != nil {
			return args, err
		}
		args = append(args, arg)
	}
}

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	mapping map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) error {
	if _, ok := c.mapping[name]; ok {
		return fmt.Errorf("command %v already registered", name)
	}
	c.mapping[name] = f
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.mapping[cmd.name]
	if !ok {
		return fmt.Errorf("unknown command %v", cmd.name)
	}
	return f(s, cmd)
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchRssFeed(ctx context.Context, url string) (*RSSFeed, error) {
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var feed RSSFeed
	err = xml.NewDecoder(resp.Body).Decode(&feed)
	if err != nil {
		return nil, err
	}

	// Unescape HTML entities
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Items {
		feed.Channel.Items[i].Title = html.UnescapeString(feed.Channel.Items[i].Title)
		feed.Channel.Items[i].Description = html.UnescapeString(feed.Channel.Items[i].Description)
	}

	return &feed, nil
}
