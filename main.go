package main

import (
	"fmt"
	"log"
	config "github.com/NZO-GB/Gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %s", err)
		return
	}


	
	err = cfg.SetUser("NZO")
	if err != nil {
		log.Fatalf("error setting user: %s", err)
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %s", err)
		return
	}
	fmt.Println(cfg)
}
