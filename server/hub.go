package main

import (
	"fmt"

	"github.com/rayhaanbhikha/file_share/packages/commands"
)

type Hub struct {
	file             *FileWrapper
	incomingCommands chan *Command
	dataAddress      string
	standardAddress  string
}

func NewHub(file *FileWrapper, standardAddress, dataAddress string) *Hub {
	return &Hub{
		incomingCommands: make(chan *Command),
		file:             file,
		dataAddress:      dataAddress,
		standardAddress:  standardAddress,
	}
}

func (h *Hub) Run() {
	for incomingCommand := range h.incomingCommands {
		h.handleCommand(incomingCommand)
	}
}

func (h *Hub) handleCommand(command *Command) {
	var message string
	fmt.Println("Command", command.id, command.body)

	switch command.id {
	case commands.DlFile:
		message = h.createFileMessage(command.body)
		break
	case commands.Ok:
		message = "OK\n"
		break
	default:
		message = "ERROR\n"
	}

	bMessage := []byte(message)

	command.sender.con.Write(bMessage)
}

func (h *Hub) createFileMessage(fileRequested string) string {
	if fileRequested != h.file.stat.Name() {
		return "ERROR: file does not exist\n"
	}

	return fmt.Sprintf("%s %s\n", "D_ADDRESS", h.dataAddress)
}

func (h *Hub) Close() {
	close(h.incomingCommands)
}
