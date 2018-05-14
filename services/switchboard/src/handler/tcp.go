package handler

import (
	"fmt"
	"log"
	"net"
)

type TCPConf struct {
	Address string
}

type TCP struct {
	conf TCPConf
	conn net.Conn
}

func NewTCP(conf TCPConf) *TCP {
	conn, err := net.Dial("tcp", conf.Address)
	if err != nil {
		log.Panicf("Unable to open TCP connection %s: %v", conf.Address, err)
	}
	return &TCP{conf, conn}
}

func (handler *TCP) Handle(err error) {
	fmt.Fprintln(handler.conn, err)
}

func (handler *TCP) Close() error {
	return handler.conn.Close()
}
