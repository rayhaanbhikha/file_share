package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rayhaanbhikha/file_share/packages/client"
	"github.com/rayhaanbhikha/file_share/packages/server"
)

var filePath string
var fileRequested string
var mode string
var address string

func init() {
	flag.StringVar(&mode, "mode", "server", "Mode can be 'server' or 'client' ")

	flag.StringVar(&filePath, "filePath", "", "Path to file")
	flag.StringVar(&address, "address", "", "Address of FTP server <ip>:<port>")

	flag.Parse()

	if mode != "server" && mode != "client" {
		fmt.Println("Mode value is incorrect")
		os.Exit(1)
	}

	if mode == "server" && filePath == "" {
		fmt.Println("filePath argument must be provided when in server mode")
		os.Exit(1)
	}

	// TODO: ip address validation.
	if mode == "client" && address == "" {
		fmt.Println("address argument must be provided when in server mode")
		os.Exit(1)
	}

}

func main() {
	if mode == "server" {
		server.Run(filePath)
	} else {
		client.Run(address)
	}

}
