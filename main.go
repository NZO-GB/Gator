package main

import (
	"log"
	"os"
	"database/sql"
	database "github.com/NZO-GB/Gator/internal/database"
	_		"github.com/lib/pq"
	config	"github.com/NZO-GB/Gator/internal/config"
)

const dbURL = "postgres://postgres:postgres@localhost:5432/gator"

func newState() (*state, error) {
	cfg, err := config.Read()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	return &state{
		cfg:	&cfg,
		db:		database.New(db),
	}, nil
}

func registerCommands() commands {
	cmds := commands {
		list: make(map[string]func(*state, command) error),
	}

	must(cmds.register("login", handlerLogin))
	must(cmds.register("register", handlerRegister))
	must(cmds.register("reset", handlerReset))
	must(cmds.register("users", handlerGetUsers))
	must(cmds.register("agg", handlerFeed))
	must(cmds.register("addfeed", handlerAddFeed))
	must(cmds.register("feeds", handlerFeeds))

	return cmds
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	stateInstance, err := newState()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Error: not enough commands")
	}

	cmds := registerCommands()

	nameString := os.Args[1]
	argumentsString := os.Args[2:]

	var commandFull command = command {
		name: nameString,
		arguments: argumentsString,
	}

	err = cmds.run(stateInstance, commandFull)
	if err != nil {
		log.Fatal("Error: ", err)
	}

}
