package terminus

import (
	"github.com/redbooth/switchboard/src/header"
	"fmt"
	"io"
	"log"
	"os/exec"
)

type ExecConf struct {
	Command string
	Args    []string
}

type Exec struct {
	conf   ExecConf
	errors chan<- error
	cmd    *exec.Cmd
	stdin  io.WriteCloser
}

func NewExec(conf ExecConf, errors chan<- error) *Exec {
	cmd := exec.Command(conf.Command, conf.Args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Panicf("Unable to open stdin for command %s: %v\n", conf.Command, err)
	}

	err = cmd.Start()
	if err != nil {
		log.Panicf("Unable to start command %s: %v\n", conf.Command, err)
	}

	return &Exec{conf, errors, cmd, stdin}
}

func (terminus *Exec) Terminate(h header.Header) {
	_, err := fmt.Fprintln(terminus.stdin, h.String())
	if err != nil {
		terminus.errors <- err
	}
}
