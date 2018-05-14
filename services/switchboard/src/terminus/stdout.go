package terminus

import (
	"github.com/redbooth/switchboard/src/header"
	"fmt"
	"os"
)

type StdoutConf struct{}

type Stdout struct {
	conf   StdoutConf
	errors chan<- error
}

func NewStdout(conf StdoutConf, errors chan<- error) *Stdout {
	return &Stdout{conf, errors}
}

func (terminus *Stdout) Terminate(h header.Header) {
	_, err := fmt.Fprintln(os.Stdout, h.String())
	if err != nil {
		terminus.errors <- err
	}
}
