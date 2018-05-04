package input

import (
	"io"
	"os"
)

type StdinConf struct{}

type Stdin struct {
	conf    StdinConf
	errors  chan<- error
	readers chan<- io.ReadCloser
}

func NewStdin(conf StdinConf, errors chan<- error, readers chan<- io.ReadCloser) *Stdin {
	return &Stdin{conf, errors, readers}
}

func (input *Stdin) Read() {
	input.readers <- os.Stdin
}
