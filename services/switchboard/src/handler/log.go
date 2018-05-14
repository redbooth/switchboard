package handler

import "log"

type LogConf struct{}

type Log struct {
	conf LogConf
}

func NewLog(conf LogConf) *Log {
	return &Log{conf}
}

func (handler *Log) Handle(err error) {
	log.Println(err)
}

func (handler *Log) Close() error {
	return nil
}
