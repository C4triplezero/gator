package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/C4triplezero/gator/internal/config"
	"github.com/C4triplezero/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	DB     *database.Queries
	Config *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := &state{
		DB:     dbQueries,
		Config: &cfg,
	}

	cmds := commands{
		Commands: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregation)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Usage: cli <command> [args...]")
	}

	cmd := command{
		Name: args[1],
		Args: args[2:],
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
