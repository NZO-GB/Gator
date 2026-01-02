package main

import (
	"fmt"
	"github.com/NZO-GB/gator/internal/config"
)

main() {
	cfg := Read()
	cfg.SetUser("NZO")
	fmt.Println(cfg)
}
