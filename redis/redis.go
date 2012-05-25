package redis

import (
	"errors"
	"github.com/simonz05/godis/redis"
)

type Config struct {
	NetAddress string `json:"netaddr"`
	DBid       int    `json:"dbid"`
	Password   string `json:"password"`
}

type Redis struct {
	client *redis.Client
	config Config
}

func New(config Config) *Redis {
	this := new(Redis)
	this.config = config
	this.client = this.newClient()
	return this
}

func (this Redis) newClient() *redis.Client {
	return redis.New(
		this.config.NetAddress,
		this.config.DBid,
		this.config.Password,
	)
}

func (this Redis) Test() error {
	test, err := this.client.Echo("test")
	if err != nil || test.String() != "test" {
		return errors.New("Please run Redis before starting")
	}
	return nil
}

func (this Redis) Key(key string) Key {
	return newKey(this, key)
}

func (this Redis) String(key string) String {
	return newString(this, key)
}

func (this Redis) Counter(key string) Counter {
	return newCounter(this, key)
}

func (this Redis) Set(key string) Set {
	return newSet(this, key)
}

func (this Redis) List(key string) List {
	return newList(this, key)
}

func (this Redis) Mutex(key string) Mutex {
	return newMutex(this, key, 1)
}

func (this Redis) Semaphore(key string, count int) Mutex {
	return newMutex(this, key, count)
}

func (this Redis) ReadWriteMutex(key string, readers int) *ReadWriteMutex {
	return newRWMutex(this, key, readers)
}

func (this Redis) Channel(key string) Channel {
	return newChannel(this, key)
}

func checkForError(err error) {
	if err != nil {
		panic(err)
	}
}

func isTimeout(err error) bool {
	return err != nil && err.Error() == "timeout expired"
}
