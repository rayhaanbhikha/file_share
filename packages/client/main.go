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
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func Run(address string) {
	con, err := net.Dial("tcp", address)
	handleErr(err)
	defer con.Close()

	// Request fileName
	con.Write([]byte("FILE_NAME\n"))
	response, err := bufio.NewReader(con).ReadBytes('\n')
	handleErr(err)

	fileRequested := string(bytes.TrimSpace(bytes.Split(response, []byte(" "))[0]))
	fmt.Println(fileRequested)

	// Returns address of DataTCP socket.
	message := fmt.Sprintf("%s %s\n", "DL_FILE", fileRequested)
	con.Write([]byte(message))
	response, err = bufio.NewReader(con).ReadBytes('\n')
	handleErr(err)

	dataAddress := string(bytes.TrimSpace(bytes.Split(response, []byte(" "))[1]))

	err = downloadFile(dataAddress, fileRequested)
	handleErr(err)
}

func downloadFile(address string, output string) error {
	con, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer con.Close()

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path.Join(cwd, outputDirectory)); os.IsNotExist(err) {
		err := os.Mkdir(path.Join(cwd, outputDirectory), 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(path.Join(outputDirectory, output), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	con.Write([]byte("STREAM\n"))
	written, err := io.Copy(file, con)
	if err != nil {
		return err
	}

	fmt.Println("Written: ", written)
	return nil
}
