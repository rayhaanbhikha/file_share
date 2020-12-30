package server

import (
	"errors"

	"github.com/rayhaanbhikha/file_share/packages/commands"
)

type Command struct {
	id     commands.CommandID
	sender *Client
	body   string
}

func NewCommand(client *Client, command string, body string) (*Command, error) {
	commandID, ok := commands.CommandMap[command]
	if !ok {
		return nil, errors.New("Command does not exist")
	}

	//FIXME: assuming body is correct always
	return &Command{sender: client, body: body, id: commandID}, nil
}
