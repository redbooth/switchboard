package filter

import "github.com/redbooth/switchboard/src/header"

type AllowConf struct{}

type Allow struct {
	conf   AllowConf
	errors chan<- error
}

func NewAllow(conf AllowConf, errors chan<- error) *Allow {
	return &Allow{conf, errors}
}

func (filter *Allow) Filter(h header.Header) bool {
	return true
}
