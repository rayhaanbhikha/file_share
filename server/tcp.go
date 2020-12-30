package main

import (
	"fmt"
	"net"

	"github.com/rayhaanbhikha/file_share/packages/utils"
)

type TCPConnection struct {
	host  string
	port  string
	label string
}

func NewTCPConnection(host, port, label string) *TCPConnection {
	return &TCPConnection{host: host, port: port, label: label}
}

func (tcp *TCPConnection) GetExposedNetworkAddess() string {
	ip, err := utils.GetIPv4Address()
	if err != nil {
		panic(err)
	}
	address := net.JoinHostPort(ip.String(), tcp.port)
	return address
}

func (tcp *TCPConnection) Run(connectionHandler func(con net.Conn)) {
	address := net.JoinHostPort(tcp.host, tcp.port)
	ln, err := net.Listen("tcp", address)
	handleErr(err)
	defer ln.Close()

	message := fmt.Sprintf("%s TCP server started on %s\n", tcp.label, address)
	fmt.Println(message)

	for {
		con, err := ln.Accept()
		handleErr(err)

		go connectionHandler(con)
	}
}
