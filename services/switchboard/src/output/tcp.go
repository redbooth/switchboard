package output

import (
	"../header"
	"log"
	"net"
)

type TCPConf struct {
	Address string
}

type TCP struct {
	conf   TCPConf
	errors chan<- error
	header header.Header
	conn   net.Conn
}

func NewTCP(conf TCPConf, errors chan<- error, h header.Header) *TCP {
	// connect to upstream
	conn, err := net.Dial("tcp", conf.Address)
	if err != nil {
		log.Panicf("Unable to open TCP connection %s: %v", conf.Address, err)
	}
	// inject header
	bytes, err := header.Bytes(h)
	if err != nil {
		log.Panicf("Unable to convert header into bytes: %v\n", err)
	}
	conn.Write(bytes)
	return &TCP{conf, errors, h, conn}
}

func (writer *TCP) Write(b []byte) (n int, err error) {
	return writer.conn.Write(b)
}

func (writer *TCP) Close() error {
	return writer.conn.Close()
}
