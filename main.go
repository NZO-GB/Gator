package main

import (
	"fmt"
	config "github.com/NZO-GB/Gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg.SetUser("NZO")
	fmt.Println(cfg)
}
