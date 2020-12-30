package server

import (
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

func Run(filePath string) {
	// Validate file.
	file := NewFile(filePath)
	err := file.Init()
	if err != nil {
		file.HandleError(err)
		return
	}
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
