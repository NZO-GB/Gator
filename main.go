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

func main() {
	var stateInstance state

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	stateInstance.cfg = &cfg

	db, err := sql.Open("postgres", dbURL)

	dbQueries := database.New(db)

	stateInstance.db = dbQueries

	if len(os.Args) < 2 {
		log.Fatal("Error: not enough commands")
	}

	cmds := commands {
		list: make(map[string]func(*state, command) error),
	}

	if err = cmds.register("login", handlerLogin); err != nil {
		log.Fatal("Register error: ", err)
	}

	if err = cmds.register("protocol", addConnectionString); err != nil {
		log.Fatal("Register error: ", err)
	}

	nameString := os.Args[1]
	argumentsString := os.Args[2:]

	var commandFull command = command {
		name: nameString,
		arguments: argumentsString,
	}

	err = cmds.run(&stateInstance, commandFull)
	if err != nil {
		log.Fatal("Error: ", err)
	}

}
