package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path"
)

const outputDirectory = "./data"

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Run(address string) {
	con, err := net.Dial("tcp", address)
	handleErr(err)
	defer con.Close()

	con.Write([]byte("FILE_NAME\n"))
	response, err := bufio.NewReader(con).ReadBytes('\n')
	fileRequested := string(bytes.TrimSpace(bytes.Split(response, []byte(" "))[0]))
	fmt.Println(fileRequested)

	// TODO: fixme.
	message := fmt.Sprintf("%s %s\n", "DL_FILE", fileRequested)
	con.Write([]byte(message))

	response, err = bufio.NewReader(con).ReadBytes('\n')
	if err != nil {
		fmt.Println("downloading error", err)
		return
	}

	dataAddress := string(bytes.TrimSpace(bytes.Split(response, []byte(" "))[1]))

	err = downloadFile(dataAddress, fileRequested)
	if err != nil {
		fmt.Println(err)
	}
}

func downloadFile(address string, output string) error {
	con, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer con.Close()

	con.Write([]byte("STREAM\n"))

	file, err := os.OpenFile(path.Join(outputDirectory, output), os.O_RDWR|os.O_CREATE, 0755)
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
