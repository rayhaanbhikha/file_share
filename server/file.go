package main

import (
	"fmt"
	"os"
)

type FileWrapper struct {
	filePath string
	stat     os.FileInfo
}

func NewFile(filePath string) *FileWrapper {
	return &FileWrapper{filePath: filePath}
}

func (f *FileWrapper) Init() error {
	stat, err := os.Stat(f.filePath)
	if err != nil {
		return err
	}
	f.stat = stat
	return nil
}

func (f *FileWrapper) HandleError(err error) {
	if os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("File %s does not exist.", f.filePath))
	} else {
		fmt.Println(err.Error())
	}
}
