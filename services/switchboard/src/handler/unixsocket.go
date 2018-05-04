package handler

import (
	"fmt"
	"log"
	"net"
)

type UnixSocketConf struct {
	Address string
}

type UnixSocket struct {
	conf UnixSocketConf
	conn net.Conn
}

func NewUnixSocket(conf UnixSocketConf) *UnixSocket {
	conn, err := net.Dial("unix", conf.Address)
	if err != nil {
		log.Panicf("Unable to open unix socket connection to %s: %v", conf.Address, err)
	}
	return &UnixSocket{conf, conn}
}

func (handler *UnixSocket) Handle(err error) {
	fmt.Fprintln(handler.conn, err)
}

func (handler *UnixSocket) Close() error {
	return handler.conn.Close()
}
