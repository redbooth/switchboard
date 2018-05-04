package input

import (
	"fmt"
	"io"
	"net"
)

type TCPConf struct {
	Port uint16
}

type TCP struct {
	conf    TCPConf
	errors  chan<- error
	readers chan<- io.ReadCloser
}

func NewTCP(conf TCPConf, errors chan<- error, readers chan<- io.ReadCloser) *TCP {
	return &TCP{conf, errors, readers}
}

func (input *TCP) Read() {
	port := fmt.Sprintf(":%d", input.conf.Port)
	if listener, err := net.Listen("tcp", port); err != nil {
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
