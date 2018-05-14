package output

import (
	"github.com/redbooth/switchboard/src/header"
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
	header header.Header
	cmd    *exec.Cmd
	stdin  io.WriteCloser
}

func NewExec(conf ExecConf, errors chan<- error, h header.Header) *Exec {
	cmd := exec.Command(conf.Command, conf.Args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Panicf("Unable to open stdin for command %s: %v\n", conf.Command, err)
	}

	err = cmd.Start()
	if err != nil {
		log.Panicf("Unable to start command %s: %v\n", conf.Command, err)
	}

	return &Exec{conf, errors, h, cmd, stdin}
}

func (writer *Exec) Write(b []byte) (n int, err error) {
	return writer.stdin.Write(b)
}

func (writer *Exec) Close() error {
	return writer.stdin.Close()
}
