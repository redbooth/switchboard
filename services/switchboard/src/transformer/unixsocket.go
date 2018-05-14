package transformer

import (
	"io"
	"net"
)

type UnixSocketConf struct {
	Address string
}

type UnixSocket struct {
	conf   UnixSocketConf
	errors chan<- error
}

func NewUnixSocket(conf UnixSocketConf, errors chan<- error) *UnixSocket {
	return &UnixSocket{conf, errors}
}

func (transformer *UnixSocket) Transform(input io.Reader) io.Reader {
	if conn, err := net.Dial("unix", transformer.conf.Address); err != nil {
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
