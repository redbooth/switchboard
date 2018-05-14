package transformer

import (
	"io"
)

type IdentityConf struct{}

type Identity struct {
	conf   IdentityConf
	errors chan<- error
}

func NewIdentity(conf IdentityConf, errors chan<- error) *Identity {
	return &Identity{conf, errors}
}

func (transformer *Identity) Transform(reader io.Reader) io.Reader {
	return reader
}
