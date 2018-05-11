package output

import (
	"github.com/redbooth/switchboard/src/header"
	"os"
)

type StdoutConf struct{}

type Stdout struct {
	conf   StdoutConf
	errors chan<- error
	header header.Header
}

func NewStdout(conf StdoutConf, errors chan<- error, h header.Header) *Stdout {
	return &Stdout{conf, errors, h}
}

func (output *Stdout) Write(b []byte) (n int, err error) {
	return os.Stdout.Write(b)
}

func (output *Stdout) Close() error {
	return nil
}
