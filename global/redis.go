package global

import (
	"github.com/HairyMezican/SimpleRedis/redis"
)

var Redis *redis.Client

func init() {
	Redis = redis.New(redis.DefaultConfiguration())
}
