package terminus

import (
	"github.com/redbooth/switchboard/src/header"
	"github.com/go-redis/redis"
)

type RedisConf struct {
	Address  string
	Password string
	Database int
	Channel  string
}

type Redis struct {
	conf   RedisConf
	errors chan<- error
	client *redis.Client
}

func NewRedis(conf RedisConf, errors chan<- error) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Address,
		Password: conf.Password,
		DB:       conf.Database,
	})
	return &Redis{conf, errors, client}
}

func (terminus *Redis) Terminate(h header.Header) {
	err := terminus.client.Publish(terminus.conf.Channel, h.String()).Err()
	if err != nil {
		terminus.errors <- err
	}
}
