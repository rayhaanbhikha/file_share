package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

const address = "127.0.0.1:8080"

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	listener, err := net.Listen("tcp", address)
	handleErr(err)
	defer listener.Close()

	fmt.Println("TCP server started")
	for {
		con, err := listener.Accept()
		handleErr(err)

		// TODO: con.read()
		go handleConnection(con)
	}
}

func handleConnection(con net.Conn) {
	defer con.Close()
	fmt.Print("---- CLIENT CONNECTED ----")
	for {
		command, err := bufio.NewReader(con).ReadString('\n')

		if err != nil {
			if err == io.EOF {
				fmt.Println("connection terminated")
				return
			} else {
				fmt.Println("ERROR: ", err)
				return
			}
		}

		switch command {
		case "DOWNLOAD\r\n":
			fmt.Println(command)
			err = streamFile(con)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("complete")
				// fmt.Fprintf(con, "\r\nCOMPLETE\r\n")
			}
			return
		case "\r\n":
			break
		default:
			message := "Command does not exist.\r\n"
			fmt.Fprint(con, message)
		}
	}
}

func streamFile(con net.Conn) error {
	file, err := os.Open("./slowup.mp3")
	// file, err := os.Open("./data.txt")
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
