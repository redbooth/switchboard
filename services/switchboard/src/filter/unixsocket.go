package filter

import (
	"github.com/redbooth/switchboard/src/header"
	"bytes"
	"encoding/binary"
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

func (filter *UnixSocket) Filter(h header.Header) bool {
	buf := new(bytes.Buffer)
	resp := make([]byte, 1)
	if conn, err := net.Dial("unix", filter.conf.Address); err != nil {
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
