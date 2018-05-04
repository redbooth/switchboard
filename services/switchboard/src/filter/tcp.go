package filter

import (
	"../header"
	"bytes"
	"encoding/binary"
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

func (filter *TCP) Filter(h header.Header) bool {
	buf := new(bytes.Buffer)
	resp := make([]byte, 1)
	if conn, err := net.Dial("tcp", filter.conf.Address); err != nil {
		filter.errors <- err
		return false
	} else if err := binary.Write(buf, binary.LittleEndian, h); err != nil {
		filter.errors <- err
		return false
	} else if _, err := conn.Write(buf.Bytes()); err != nil {
		filter.errors <- err
		return false
	} else if _, err := conn.Read(resp); err != nil {
		filter.errors <- err
		return false
	} else {
		return resp[0] == 0x00
	}
}
