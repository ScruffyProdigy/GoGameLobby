package redis

import (
	"../loadconfiguration"
)

var Client *Redis

func init() {
	var c Config
	configurations.Load("redis", &c)
	Client = New(c)
}
