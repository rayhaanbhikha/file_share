package main

import (
	"fmt"
	"net"
)

type TCPConnection struct {
	address, label string
}

func NewTCPConnection(address, label string) *TCPConnection {
	return &TCPConnection{address: address, label: label}
}

func (tcp *TCPConnection) run(connectionHandler func(con net.Conn)) {
	ln, err := net.Listen("tcp", tcp.address)
	handleErr(err)
	defer ln.Close()

	message := fmt.Sprintf("%s TCP server started on %s\n", tcp.label, ln.Addr().String())
	fmt.Println(message)

	for {
		con, err := ln.Accept()
		handleErr(err)

		go connectionHandler(con)
	}
}
