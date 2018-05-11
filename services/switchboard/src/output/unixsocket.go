package output

import (
	"github.com/redbooth/switchboard/src/header"
	"log"
	"net"
)

type UnixSocketConf struct {
	Address string
}

type UnixSocket struct {
	conf   UnixSocketConf
	errors chan<- error
	header header.Header
	conn   net.Conn
}

func NewUnixSocket(conf UnixSocketConf, errors chan<- error, h header.Header) *UnixSocket {
	// connect to upstream
	conn, err := net.Dial("unix", conf.Address)
	if err != nil {
		log.Panicf("Unable to open unix socket connection to %s: %v", conf.Address, err)
	}
	// inject header
	bytes, err := header.Bytes(h)
	if err != nil {
		log.Panicf("Unable to convert header into bytes: %v\n", err)
	}
	conn.Write(bytes)
	return &UnixSocket{conf, errors, h, conn}
}

func (writer *UnixSocket) Write(b []byte) (n int, err error) {
	return writer.conn.Write(b)
}

func (writer *UnixSocket) Close() error {
	return writer.conn.Close()
}
