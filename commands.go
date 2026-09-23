package main

import (
	"fmt"

	"github.com/Solaire89/gator/internal/config"
	"github.com/Solaire89/gator/internal/database"
)

// Shared toolbox that are used by each handler
type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	Name string
	Args []string
}

// This struct is like a receptionist that takes a request (the command) and
// sends that request to the appropriate function
type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("error: command not found")
	}
	return handler(s, cmd)
}
