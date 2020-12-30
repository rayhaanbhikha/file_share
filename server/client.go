package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

type Client struct {
	con      net.Conn
	outbound chan<- *Command
}

func NewClient(con net.Conn, outbound chan<- *Command) *Client {
	return &Client{con: con, outbound: outbound}
}

func (c *Client) read() {
	fmt.Println("Client address: ", c.con.RemoteAddr())
	for {
		msg, err := bufio.NewReader(c.con).ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection terminated")
				return
			}

			fmt.Println("ERROR: ", err)
			// TODO: should you terminate connection if there was an error?
		}
		go c.handleMessage(msg)
	}
}

func (c *Client) handleMessage(message []byte) {
	cmd := bytes.ToUpper(bytes.TrimSpace(bytes.Split(message, []byte(" "))[0]))
	args := bytes.TrimSpace(bytes.TrimPrefix(message, cmd))

	// TODO: implement proper cmd and args validation.

	command := string(cmd)
	arguments := string(args)
	fmt.Println("COMMAND: ", command, arguments)

	formattedCommand, err := NewCommand(c, command, arguments)

	if err != nil {
		// TODO: handle errors better.
		c.con.Write([]byte("ERROR\n"))
		return
	}

	c.outbound <- formattedCommand
}
