package main

import (
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
	con, err := net.Dial("tcp", "127.0.0.1:8080")
	handleErr(err)

	file, err := os.OpenFile("data.mp3", os.O_RDWR|os.O_CREATE, 0755)

	downloadMessage := []byte("DOWNLOAD\r\n")

	con.Write(downloadMessage)

	written, err := io.Copy(file, con)
	handleErr(err)

	fmt.Println("Written: ", written)
}
