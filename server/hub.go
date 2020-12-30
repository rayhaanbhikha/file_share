package main

import (
	"fmt"
	"io"
	"net"

	"github.com/rayhaanbhikha/file_share/packages/commands"
	"github.com/rayhaanbhikha/file_share/packages/utils"
)

type Hub struct {
	file             *FileWrapper
	incomingCommands chan *Command
	dataAddress      string
	standardAddress  string
}

func NewHub(file *FileWrapper, ip, standardPort, dataPort string) *Hub {
	return &Hub{
		incomingCommands: make(chan *Command),
		file:             file,
		standardAddress:  net.JoinHostPort(ip, standardPort),
		dataAddress:      net.JoinHostPort(ip, dataPort),
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
		h.respondToFileDownload(command)
		return
	case commands.Ok:
		message = "OK"
		command.sender.con.Write(utils.FormatStreamMessage(message))
		break
	case commands.Stream:
		h.streamFile(command)
		return
	default:
		message = "ERROR"
		command.sender.con.Write(utils.FormatStreamMessage(message))
		return
	}
}

func (h *Hub) respondToFileDownload(command *Command) {
	fileRequested := command.body
	var message string

	if fileRequested != h.file.stat.Name() {
		message = "ERROR: file does not exist"
	} else {
		message = fmt.Sprintf("%s %s", "D_ADDRESS", h.dataAddress)
	}

	command.sender.con.Write(utils.FormatStreamMessage(message))
}

func (h *Hub) streamFile(command *Command) {
	var message string

	written, err := io.Copy(command.sender.con, h.file.file)
	if err != nil {
		fmt.Println("STREAMING error: ", err.Error())
		message = "ERROR"
		command.sender.con.Write(utils.FormatStreamMessage(message))
	} else {
		fmt.Println("Written: ", written)
	}

	command.sender.closeConnection()
}

func (h *Hub) Close() {
	close(h.incomingCommands)
}
