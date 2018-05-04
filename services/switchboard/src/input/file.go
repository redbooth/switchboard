package input

import (
	"io"
	"os"
)

type FileConf struct {
	Name string
}

type File struct {
	conf    FileConf
	errors  chan<- error
	readers chan<- io.ReadCloser
}

func NewFile(conf FileConf, errors chan<- error, readers chan<- io.ReadCloser) *File {
	return &File{conf, errors, readers}
}

func (input *File) Read() {
	if file, err := os.Open(input.conf.Name); err != nil {
		input.errors <- err
	} else {
		input.readers <- file
	}
}
