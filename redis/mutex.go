package redis

import (
	"../mutex"
	"github.com/simonz05/godis/redis"
)

type redismutex struct {
	client *redis.Client
	key    string
}

func Mutex(key string) (mutex.Mutex, error) {
	return Semaphore(key, 1)
}

func Semaphore(key string, count int) (mutex.Mutex, error) {
	oldvalue, err := Client.Getset(key+":Init", "initialized")
	if err != nil {
		return nil, err
	}

	if oldvalue.String() != "initialized" {
		for i := 0; i < count; i++ {
			_, err := Client.Rpush(key, i+1)
			if err != nil {
				Client.Del(key + ":Init")
				Client.Del(key)
				return nil, err
			}
		}
	}

	m := new(redismutex)
	m.client = Client
	m.key = key
	return m, nil
}

func (this *redismutex) Try(action func()) bool {
	old, err := this.client.Lpop(this.key)
	if err != nil {
		panic(err)
	}
	if old.Int64() == 0 {
		return false
	}

	defer func() {
		_, err := this.client.Rpush(this.key, old.Int64())
		if err != nil {
			panic(err)
		}
	}()

	action()

	return true
}

func (this *redismutex) Force(action func()) {
	old, err := this.client.Blpop([]string{this.key}, 0)
	if err != nil {
		panic(err)
	}

	defer func() {
		_, err := this.client.Rpush(this.key, old)
		if err != nil {
			panic(err)
		}
	}()

	action()
}
