package redis

import (
	"../loadconfiguration"
	"github.com/simonz05/godis/redis"
)

var Client *redis.Client

type redisConfig struct {
	NetAddress string `json:"netaddr"`
	DBid       int    `json:"dbid"`
	Password   string `json:"password"`
}

func (r redisConfig) Redis() *redis.Client {
	return redis.New(
		r.NetAddress,
		r.DBid,
		r.Password,
	)
}

func init() {
	var config redisConfig

	configurations.Load("redis", &config)
	Client = config.Redis()

	test, err := Client.Echo("test")
	if err != nil || test.String() != "test" {
		panic("Please run Redis before executing this")
	}
}
