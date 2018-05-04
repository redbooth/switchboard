package handler

import (
	"fmt"
	"os"
)

type StderrConf struct{}

type Stderr struct {
	conf StderrConf
}

func NewStderr(conf StderrConf) *Stderr {
	return &Stderr{conf}
}

func (handler *Stderr) Handle(err error) {
	fmt.Fprintln(os.Stderr, err)
}

func (output *Stderr) Close() error {
	return nil
}
