package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const address = "0.0.0.0:8080"
const dataAddress = "0.0.0.0:8081"

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	hub := NewHub()
	go hub.Run()
	defer hub.Close()

	listener, err := net.Listen("tcp", address)
	handleErr(err)
	defer listener.Close()
	fmt.Println("TCP server started")

	go dataTCPConnection()

	for {
		con, err := listener.Accept()
		handleErr(err)
		fmt.Println("DOWNLOADING DATA")
		client := NewClient(con, hub.incomingCommands)
		client.read()
	}
}

func dataTCPConnection() {
	listener, err := net.Listen("tcp", dataAddress)
	handleErr(err)
	defer listener.Close()
	fmt.Println("TCP DATA server started")

	for {
		con, err := listener.Accept()
		handleErr(err)

		err = streamFile(con)
		if err != nil {
			fmt.Println("ERROR: ", err)
			con.Write([]byte("ERROR\n"))
		}
		fmt.Println("closing connection")
		con.Close()
	}
}

func streamFile(con net.Conn) error {
	// TODO: check said file exists.
	file, err := os.Open("./1gbfile")
	handleErr(err)
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	fmt.Println("FILE INFO: ")
	fmt.Println("Name: ", stat.Name())
	fmt.Println("Size: ", stat.Size())
	written, err := io.Copy(con, file)
	if err != nil {
		return err
	}
	fmt.Println("Written: ", written)
	return nil
}
