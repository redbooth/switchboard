package handler

import (
	"fmt"
	"log"
	"os"
)

type FileConf struct {
	Name string
}

type File struct {
	conf FileConf
	file *os.File
}

func NewFile(conf FileConf) *File {
	filename := conf.Name
	file, err := os.Create(filename)
	if err != nil {
		log.Panicf("Unable to open file %s: %v\n", filename, err)
	}
	return &File{conf, file}
}

func (handler *File) Handle(err error) {
	fmt.Fprintln(handler.file, err)
}

func (handler *File) Close() error {
	return handler.file.Close()
}
