package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
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

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	standardTCPConnection := NewTCPConnection(address, "Standard")
	dataTCPConnection := NewTCPConnection(dataAddress, "Data")

	go standardTCPConnection.run(func(con net.Conn) {
		client := NewClient(con, hub.incomingCommands)
		client.read()
	})

	go dataTCPConnection.run(func(con net.Conn) {
		defer con.Close()
		fmt.Println("DOWNLOADING DATA")
		err := streamFile(con)
		if err != nil {
			fmt.Println("ERROR: ", err)
			con.Write([]byte("ERROR\n"))
		}
		fmt.Println("closing connection")
	})

	s := <-signalChannel
	fmt.Println(s)
}

func streamFile(con net.Conn) error {
	// TODO: check said file exists.
	// file, err := os.Open("./1gbfile")
	file, err := os.Open("./slowup.mp3")
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
