package main

import (
	"../../../src/filter"
	"../../../src/header"
	"log"
)

type AuthConf struct {
	RedisAddress string
	RedisDatabase uint8
}

type Auth struct {
	conf AuthConf
}

func Conf() AuthConf {
	return AuthConf{}
}

func Constructor(conf filter.Conf) filter.Filter {
	c, ok := conf.(AuthConf)
	if !ok {
		log.Panicf("Unexpected type for Configuration %v: %T", conf, conf)
	}
	return &Auth{c}
}

func (filter *Auth) Filter(h header.Header, errors chan<- error) bool {
	// TODO: verify that user is logged in and has access to the convo
	return true
}
