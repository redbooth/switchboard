package filter

import "github.com/redbooth/switchboard/src/header"

type DenyConf struct{}

type Deny struct {
	conf   DenyConf
	errors chan<- error
}

func NewDeny(conf DenyConf, errors chan<- error) *Deny {
	return &Deny{conf, errors}
}

func (filter *Deny) Filter(h header.Header) bool {
	return false
}
