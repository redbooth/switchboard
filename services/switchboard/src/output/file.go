package output

import (
	"../header"
	"log"
	"os"
	"path"
)

type FileConf struct {
	Directory string
	Extension string
}

type File struct {
	conf   FileConf
	errors chan<- error
	header header.Header
	file   *os.File
}

func NewFile(conf FileConf, errors chan<- error, h header.Header) *File {
	err := os.MkdirAll(conf.Directory, 0600)
	if err != nil {
		log.Panicf("Unable to create directory %s: %v\n", conf.Directory, err)
	}
	filename := path.Join(conf.Directory, h.String())
	if len(conf.Extension) > 0 {
		filename += "." + conf.Extension
	}
	file, err := os.Create(filename)
	if err != nil {
		log.Panicf("Unable to open file %s: %v\n", filename, err)
	}
	return &File{conf, errors, h, file}
}

func (output *File) Write(b []byte) (n int, err error) {
	return output.file.Write(b)
}

func (output *File) Close() error {
	return output.file.Close()
}
