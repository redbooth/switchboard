package handler

import (
	"net/http"
	"strings"
)

type HTTPConf struct {
	Address string
}

type HTTP struct {
	conf HTTPConf
}

func NewHTTP(conf HTTPConf) *HTTP {
	return &HTTP{conf}
}

func (handler *HTTP) Handle(err error) {
	http.Post(handler.conf.Address, "text/plain", strings.NewReader(err.Error()))
}

func (handler *HTTP) Close() error {
	return nil
}
