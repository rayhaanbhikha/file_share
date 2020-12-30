package server

import (
	"fmt"
	"net"
)

type TCPConnection struct {
	host  string
	port  string
	label string
}

func NewTCPConnection(host, port, label string) *TCPConnection {
	return &TCPConnection{host: host, port: port, label: label}
}

func (tcp *TCPConnection) Run(hub *Hub) {
	address := net.JoinHostPort(tcp.host, tcp.port)
	ln, err := net.Listen("tcp", address)
	handleErr(err)
	defer ln.Close()

	message := fmt.Sprintf("%s TCP server started on %s\n", tcp.label, address)
	fmt.Println(message)

	for {
		con, err := ln.Accept()
		handleErr(err)

		client := NewClient(con, hub.incomingCommands)
		go client.read()
	}
}
