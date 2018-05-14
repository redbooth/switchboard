package input

import (
	"io"
	"net"
)

type UnixSocketConf struct {
	Address string
}

type UnixSocket struct {
	conf    UnixSocketConf
	errors  chan<- error
	readers chan<- io.ReadCloser
}

func NewUnixSocket(conf UnixSocketConf, errors chan<- error, readers chan<- io.ReadCloser) *UnixSocket {
	return &UnixSocket{conf, errors, readers}
}

func (input *UnixSocket) Read() {
	if listener, err := net.Listen("unix", input.conf.Address); err != nil {
		input.errors <- err
	} else {
		for {
			if conn, err := listener.Accept(); err != nil {
				input.errors <- err
			} else {
				input.readers <- conn
			}
		}
	}
}
