package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	con, err := net.Dial("tcp", "127.0.0.1:8080")
	handleErr(err)
	defer con.Close()

	// TODO: fixme.
	downloadMessage := []byte("DL_FILE some-file\n")
	con.Write(downloadMessage)

	response, err := bufio.NewReader(con).ReadString('\n')
	if err != nil {
		fmt.Println("downloading error", err)
		return
	}

	dataAddress := strings.Split(strings.TrimRight(response, "\n"), " ")[1]
	fmt.Println(dataAddress)

	err = downloadFile(dataAddress)
	if err != nil {
		fmt.Println(err)
	}
}

func downloadFile(address string) error {
	con, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		return err
	}
	defer con.Close()

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
