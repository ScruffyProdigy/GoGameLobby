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

type Root struct {
	client *redis.Client
	config Config
}

func New(config Config) *Root {
	this := new(Root)
	this.config = config
	this.client = this.newClient()
	err := this.Test()
	checkForError(err)
	return this
}

func (this Root) newClient() *redis.Client {
	return redis.New(
		this.config.NetAddress,
		this.config.DBid,
		this.config.Password,
	)
}

func (this Root) Test() error {
	test, err := this.client.Echo("test")
	if err != nil || test.String() != "test" {
		return errors.New("Please run Redis before starting")
	}
	return nil
}

func (this Root) Key(key string) Key {
	return newKey(this, key)
}

func (this Root) String(key string) String {
	return newString(this, key)
}

func (this Root) Integer(key string) Integer {
	return newInteger(this, key)
}

func (this Root) Set(key string) Set {
	return newSet(this, key)
}

func (this Root) List(key string) List {
	return newList(this, key)
}

func (this Root) Mutex(key string) Mutex {
	return newMutex(this, key, 1)
}

func (this Root) Semaphore(key string, count int) Mutex {
	return newMutex(this, key, count)
}

func (this Root) ReadWriteMutex(key string, readers int) *ReadWriteMutex {
	return newRWMutex(this, key, readers)
}

func (this Root) Channel(key string) Channel {
	return newChannel(this, key)
}

func (this Root) Prefix(key string) Prefix {
	return newPrefix(this,key)
}

func checkForError(err error) {
	if err != nil {
		panic(err)
	}
}

func isTimeout(err error) bool {
	return err != nil && err.Error() == "timeout expired"
}


func ftoa(f float64) string {
	return strconv.FormatFloat(f,'f',-1,64)
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func atoi(s string) int64 {
	i,e := strconv.ParseInt(s,10,64)
	if e != nil {
		return 0
	}
	return i
}

func atof(s string) float64 {
	f,e := strconv.ParseFloat(s,64)
	if e != nil {
		return 0
	}
	return f
}
