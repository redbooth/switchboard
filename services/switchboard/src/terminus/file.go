package terminus

import (
	"github.com/redbooth/switchboard/src/header"
	"fmt"
	"log"
	"os"
)

type FileConf struct {
	Name string
}

type File struct {
	conf   FileConf
	errors chan<- error
	file   *os.File
}

func NewFile(conf FileConf, errors chan<- error) *File {
	filename := conf.Name
	file, err := os.Create(filename)
	if err != nil {
		log.Panicf("Unable to open file %s: %v\n", filename, err)
	}
	return &File{conf, errors, file}
}

func (terminus *File) Terminate(h header.Header) {
	_, err := fmt.Fprintln(terminus.file, h.String())
	if err != nil {
		terminus.errors <- err
	}
}
