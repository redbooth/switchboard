package terminus

import (
	"../header"
	"fmt"
	"log"
	"net"
)

type UnixSocketConf struct {
	Address string
}

type UnixSocket struct {
	conf   UnixSocketConf
	errors chan<- error
	conn   net.Conn
}

func NewUnixSocket(conf UnixSocketConf, errors chan<- error) *UnixSocket {
	conn, err := net.Dial("unix", conf.Address)
	if err != nil {
		log.Panicf("Unable to open unix socket connection %s: %v", conf.Address, err)
	}
	return &UnixSocket{conf, errors, conn}
}

func (terminus *UnixSocket) Terminate(h header.Header) {
	_, err := fmt.Fprintln(terminus.conn, h.String())
	if err != nil {
		terminus.errors <- err
	}
}
