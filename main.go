package main

import (
	"fmt"
	"log"
	"os"

	config "github.com/NZO-GB/Gator/internal/config"
)

func main() {
	var stateInstance state

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %s", err)
		return
	}

	stateInstance.cfg = &cfg

	if len(os.Args) < 2 {
		log.Fatal("Error: not enough commands")
		return
	}

	commandsInstance := commands {
		list: make(map[string]func(*state, command) error),
	}

	if err = commandsInstance.register("login", handlerLogin); err != nil {
		log.Fatal("Register error: ", err)
	}

	if err = commandsInstance.register("protocol", addConnectionString); err != nil {
		log.Fatal("Register error: ", err)
	}

	nameString := os.Args[1]
	argumentsString := os.Args[2:]

	var commandFull command = command {
		name: nameString,
		arguments: argumentsString,
	}

	fmt.Println(commandFull)

	err = commandsInstance.run(&stateInstance, commandFull)

	if err != nil {
		log.Fatal("Error: ", err)
	}

}
