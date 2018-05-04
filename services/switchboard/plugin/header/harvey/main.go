package main

import (
	"../../../src/header"
	"log"
	"encoding/hex"
)

type HarveyConf struct {}

type Harvey struct {
	conf         HarveyConf
	User         [16]byte
	Token        [16]byte
	Conversation [16]byte
}

func Conf() HarveyConf {
	return HarveyConf{}
}

func Constructor(conf header.Conf) header.Header {
	c, ok := conf.(HarveyConf)
	if !ok {
		log.Panicf("Unexpected type for Configuration %v: %T", conf, conf)
	}
	return &Harvey{conf:c}
}

func (header *Harvey) Id() []byte {
	return header.Conversation[:]
}

func (header *Harvey) String() string {
	return hex.EncodeToString(header.Id())
}
