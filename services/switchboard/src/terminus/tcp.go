package terminus

import (
	"github.com/redbooth/switchboard/src/header"
	"fmt"
	"log"
	"net"
)

type TCPConf struct {
	Address string
}

type TCP struct {
	conf   TCPConf
	errors chan<- error
	conn   net.Conn
}

func NewTCP(conf TCPConf, errors chan<- error) *TCP {
	conn, err := net.Dial("tcp", conf.Address)
	if err != nil {
		log.Panicf("Unable to open TCP connection %s: %v", conf.Address, err)
	}
	return &TCP{conf, errors, conn}
}

func (terminus *TCP) Terminate(h header.Header) {
	_, err := fmt.Fprintln(terminus.conn, h.String())
	if err != nil {
		terminus.errors <- err
	}
}
