package main

import "os"

type FileWrapper struct {
	filePath string
}

func NewFile(filePath string) *FileWrapper {
	return &FileWrapper{filePath: filePath}
}

func (f *FileWrapper) fileExists() bool {
	if _, err := os.Stat(f.filePath); os.IsNotExist(err) {
		return false
	}
	return true
}
