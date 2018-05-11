package filter

import (
	"github.com/redbooth/switchboard/src/header"
	"bytes"
	"encoding/binary"
	"net/http"
)

type HTTPConf struct {
	Address string
}

type HTTP struct {
	conf   HTTPConf
	errors chan<- error
}

func NewHTTP(conf HTTPConf, errors chan<- error) *HTTP {
	return &HTTP{conf, errors}
}

func (filter *HTTP) Filter(h header.Header) bool {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, h); err != nil {
		filter.errors <- err
		return false
	} else if resp, err := http.Post(filter.conf.Address, "application/octet-stream", buf); err != nil {
		filter.errors <- err
		return false
	} else {
		return resp.StatusCode == http.StatusOK
	}
}
