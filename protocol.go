package main

import (
	"fmt"
)

const (
	CommandGET = "GET"
	CommandSET = "SET"
)

type Command interface {
}

type SetCommand struct {
	key, val []byte
}

type GetCommand struct {
	key, val []byte
}

func parseCommand(rawMsg string) (Command, error) {
	return nil, fmt.Errorf("parseCommand error")
}
