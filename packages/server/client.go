package server

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
	done     chan int
}

func NewClient(con net.Conn, outbound chan<- *Command) *Client {
	return &Client{con: con, outbound: outbound, done: make(chan int, 0)}
}

func (c *Client) read() {
	defer c.con.Close()
	fmt.Println("Client address: ", c.con.RemoteAddr())

	go func() {
		for {
			msg, err := bufio.NewReader(c.con).ReadBytes('\n')
			if err != nil {
				if err == io.EOF || c.con == nil {
					fmt.Println("connection terminated")
					return
				}
				fmt.Println("ERROR: ", err)
			} else {
				go c.handleMessage(msg)
			}
		}
	}()

	<-c.done
	return
}

func (c *Client) handleMessage(message []byte) {
	cmd := bytes.ToUpper(bytes.TrimSpace(bytes.Split(message, []byte(" "))[0]))
	args := bytes.TrimSpace(bytes.TrimPrefix(message, cmd))

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

func (c *Client) closeConnection() {
	c.done <- 1
	c.con = nil
}
