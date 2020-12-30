package main

import "errors"

type CommandID int

const (
	dlFile CommandID = iota
	dAddress
	ok
)

type Command struct {
	id     CommandID
	sender *Client
	body   string
}

func NewCommand(client *Client, command string, body string) (*Command, error) {
	var commandID CommandID
	switch command {
	case "DL_FILE":
		commandID = dlFile
	case "D_ADDRESS":
		commandID = dAddress
	case "OK":
		commandID = ok
	default:
		return nil, errors.New("Command does not exist")
	}

	//FIXME: assuming body is correct always
	return &Command{sender: client, body: body, id: commandID}, nil
}
