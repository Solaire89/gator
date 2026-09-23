package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Solaire89/gator/internal/config"
	"github.com/Solaire89/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	// Making an instance of a cfg struct
	config, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	db, err := sql.Open("postgres", config.DbURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	st := state{
		db:  dbQueries,
		cfg: &config,
	}
	c := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}
	if err := c.run(&st, cmd); err != nil {
		log.Fatal(err)
	}
}
