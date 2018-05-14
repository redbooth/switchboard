package transformer

import (
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
}

func NewExec(conf ExecConf, errors chan<- error) *Exec {
	return &Exec{conf, errors}
}

func (transformer *Exec) Transform(input io.Reader) io.Reader {
	cmd := exec.Command(transformer.conf.Command, transformer.conf.Args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		transformer.errors <- err
		log.Panicf("Unable to open stdin for command %s: %v\n", transformer.conf.Command, err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		transformer.errors <- err
		log.Panicf("Unable to open stdout for command %s: %v\n", transformer.conf.Command, err)
	}

	go func() {
		defer stdin.Close()
		if _, err := io.Copy(stdin, input); err != nil {
			transformer.errors <- err
		}
	}()

	go func() {
		defer stdout.Close()
		if err := cmd.Run(); err != nil {
			transformer.errors <- err
		}
	}()

	return stdout
}
