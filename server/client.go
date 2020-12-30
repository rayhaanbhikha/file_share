package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
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
		message, err := bufio.NewReader(c.con).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("connection terminated")
				return
			}

			fmt.Println("ERROR: ", err)
			// TODO: should you terminate connection if there was an error?
		}
		go c.handleMessage(message)
	}
}

func (c *Client) handleMessage(message string) {
	parsedCmd := strings.Split(message, " ")

	//FIXME:
	if len(parsedCmd) != 2 {
		c.con.Write([]byte("ERROR\n"))
		return
	}

	cmd := strings.ToUpper(parsedCmd[0])
	body := parsedCmd[1]
	command, err := NewCommand(c, cmd, body)

	if err != nil {
		// TODO: handle errors better.
		c.con.Write([]byte("ERROR\n"))
		return
	}

	c.outbound <- command
}
