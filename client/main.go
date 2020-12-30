package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	fileRequested := "slowup.mp3"

	con, err := net.Dial("tcp", "127.0.0.1:8080")
	handleErr(err)
	defer con.Close()

	// TODO: fixme.
	message := fmt.Sprintf("%s %s\n", "DL_FILE", fileRequested)
	con.Write([]byte(message))

	response, err := bufio.NewReader(con).ReadBytes('\n')
	if err != nil {
		fmt.Println("downloading error", err)
		return
	}

	address := string(bytes.TrimSpace(bytes.Split(response, []byte(" "))[1]))

	err = downloadFile(address)
	if err != nil {
		fmt.Println(err)
	}
}

func downloadFile(address string) error {
	con, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer con.Close()

	con.Write([]byte("STREAM\n"))

	file, err := os.OpenFile("data.mp3", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	written, err := io.Copy(file, con)
	if err != nil {
		return err
	}

	fmt.Println("Written: ", written)
	return nil
}
