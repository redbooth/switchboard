package terminus

import (
	"../header"
	"log"
)

type LogConf struct{}

type Log struct {
	conf   LogConf
	errors chan<- error
}

func NewLog(conf LogConf, errors chan<- error) *Log {
	return &Log{conf, errors}
}

func (terminus *Log) Terminate(h header.Header) {
	log.Println(h.String())
}
