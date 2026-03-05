package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	Commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.Commands[cmd.Name]
	if !ok {
		return errors.New("error unknown command:" + cmd.Name)
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.Commands[name] = f
}
