package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/rayhaanbhikha/file_share/packages/utils"
)

const standardPort = "8080"
const dataPort = "8081"
const defaultIPAddress = "0.0.0.0"

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

var filePath string

func init() {
	flag.StringVar(&filePath, "filePath", "", "Path to file")

	flag.Parse()

	if filePath == "" {
		fmt.Println("filePath argument must be provided")
		os.Exit(1)
	}
}

func main() {
	// Validate file.
	file := NewFile(filePath)
	err := file.Init()
	if err != nil {
		file.HandleError(err)
		return
	}
	defer file.Close()
	// #############################################

	ip, err := utils.GetIPv4Address()
	if err != nil {
		panic(err)
	}

	hub := NewHub(file, ip.String(), standardPort, dataPort)
	go hub.Run()
	defer hub.Close()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)

	standardTCPConnection := NewTCPConnection(defaultIPAddress, standardPort, "Standard")
	dataTCPConnection := NewTCPConnection(defaultIPAddress, dataPort, "Data")

	go standardTCPConnection.Run(hub)
	go dataTCPConnection.Run(hub)

	s := <-signalChannel
	fmt.Println(s)
}
