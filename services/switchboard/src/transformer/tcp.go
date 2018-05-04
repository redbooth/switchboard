package transformer

import (
	"io"
	"net"
)

type TCPConf struct {
	Address string
}

type TCP struct {
	conf   TCPConf
	errors chan<- error
}

func NewTCP(conf TCPConf, errors chan<- error) *TCP {
	return &TCP{conf, errors}
}

func (transformer *TCP) Transform(input io.Reader) io.Reader {
	if conn, err := net.Dial("tcp", transformer.conf.Address); err != nil {
		transformer.errors <- err
		return nil
	} else {
		go func() {
			if _, err := io.Copy(conn, input); err != nil {
				transformer.errors <- err
			}
		}()
		return conn
	}
}
