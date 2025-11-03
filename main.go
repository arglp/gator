package main

import (
	"github.com/arglp/gator/internal/config"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	cfg.SetUser("leo")
	cfg, err = config.Read()
	fmt.Println(cfg)
}