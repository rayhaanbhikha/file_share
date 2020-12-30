package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
)

const address = "0.0.0.0"
const dataAddress = "0.0.0.0"

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	// Validate file.
	filePath := "./slowup.mp3"
	file := NewFile(filePath)
	err := file.Init()
	if err != nil {
		file.HandleError(err)
		return
	}

	// #############################################

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)

	standardTCPConnection := NewTCPConnection(address, "8080", "Standard")
	dataTCPConnection := NewTCPConnection(dataAddress, "8081", "Data")

	hub := NewHub(file, standardTCPConnection.GetExposedNetworkAddess(), dataTCPConnection.GetExposedNetworkAddess())
	go hub.Run()
	defer hub.Close()

	fmt.Println(standardTCPConnection.GetExposedNetworkAddess())

	go standardTCPConnection.Run(func(con net.Conn) {
		client := NewClient(con, hub.incomingCommands)
		client.read()
	})

	go dataTCPConnection.Run(handleDataTCPConnection)

	s := <-signalChannel
	fmt.Println(s)
}

func handleDataTCPConnection(con net.Conn) {
	defer con.Close()
	fmt.Println("DOWNLOADING DATA")
	err := streamFile(con)
	if err != nil {
		fmt.Println("ERROR: ", err)
		con.Write([]byte("ERROR\n"))
	}
	fmt.Println("closing connection")
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
