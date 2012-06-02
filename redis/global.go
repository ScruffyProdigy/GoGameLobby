package redis

import (
	"../loadconfiguration"
)

var Redis Prefix

func init() {
	var c Config
	configurations.Load("redis", &c)
	Redis = New(c)
}
