package main

import (
	"github.com/arglp/gator/internal/config"
	"github.com/arglp/gator/internal/database"
	"fmt"
	"log"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
)

type state struct {
	db 	*database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	s := state{
		cfg: &cfg,
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.db = database.New(db)

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	args := os.Args
	if len(args) < 2 {
		fmt.Errorf("argument not provided")
		os.Exit(1)
	}
	cmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}