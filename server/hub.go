package main

import "fmt"

type Hub struct {
	incomingCommands chan *Command
}

func NewHub() *Hub {
	return &Hub{incomingCommands: make(chan *Command)}
}

func (h *Hub) Run() {
	for incomingCommand := range h.incomingCommands {
		fmt.Println("command: ", incomingCommand.id)
		h.handleCommand(incomingCommand)
	}
}

func (h *Hub) handleCommand(command *Command) {
	var message string
	switch command.id {
	case 0:
		// FIXME: should be in config.
		message = "D_ADDRESS 127.0.0.1:8081\n"
		break
	case 2:
		message = "OK\n"
		break
	default:
		message = "ERROR\n"
	}

	bMessage := []byte(message)

	command.sender.con.Write(bMessage)
}

func (h *Hub) Close() {
	close(h.incomingCommands)
}
