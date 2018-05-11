package filter

import (
	"github.com/redbooth/switchboard/src/header"
	"bytes"
	"encoding/binary"
	"os/exec"
)

type ExecConf struct {
	Command string
}

type Exec struct {
	conf   ExecConf
	errors chan<- error
}

func NewExec(conf ExecConf, errors chan<- error) *Exec {
	return &Exec{conf, errors}
}

func (filter *Exec) Filter(h header.Header) bool {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, h); err != nil {
		filter.errors <- err
		return false
	} else if err := exec.Command(filter.conf.Command, buf.String()).Run(); err != nil {
		filter.errors <- err
		return false
	} else {
		return true
	}
}
