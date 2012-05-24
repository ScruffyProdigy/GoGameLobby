package redis

import (
	"../loadconfiguration"
	"errors"
	"github.com/simonz05/godis/redis"
)

var Client *redis.Client

type redisConfig struct {
	NetAddress string `json:"netaddr"`
	DBid       int    `json:"dbid"`
	Password   string `json:"password"`
}

var config redisConfig

func (r redisConfig) Redis() *redis.Client {
	return redis.New(
		r.NetAddress,
		r.DBid,
		r.Password,
	)
}

func Subscribe(url string) (*redis.Sub, error) {
	r := config.Redis()
	return r.Subscribe(url)
}

func init() {
	configurations.Load("redis", &config)
	Client = config.Redis()

	test, err := Client.Echo("test")
	if err != nil || test.String() != "test" {
		panic("Please run Redis before executing this")
	}
}

var TimeoutError = errors.New("Timeout Error")

func checkForError(err error) {
	if err != nil {
		panic(err)
	}
}

func isTimeout(err *error) bool {
	if (*err).Error() == "timeout expired" {
		*err = TimeoutError
		return true
	}
	return false
}
