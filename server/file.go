package main

import (
	"fmt"
	"os"
)

type FileWrapper struct {
	filePath string
	stat     os.FileInfo
	file     *os.File
}

func NewFile(filePath string) *FileWrapper {
	return &FileWrapper{filePath: filePath}
}

func (f *FileWrapper) Init() error {
	file, err := os.Open(f.filePath)
	handleErr(err)
	f.file = file

	stat, err := os.Stat(f.filePath)
	if err != nil {
		return err
	}
	f.stat = stat
	return nil
}

func (f *FileWrapper) Close() {
	f.file.Close()
}

func (f *FileWrapper) HandleError(err error) {
	if os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("File %s does not exist.", f.filePath))
	} else {
		fmt.Println(err.Error())
	}
}
