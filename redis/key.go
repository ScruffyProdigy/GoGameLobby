package redis

import (
	"github.com/simonz05/godis/redis"
	"time"
)

type Key struct {
	key    string
	client *redis.Client
}

func newKey(client Root, key string) Key {
	return Key{
		key:    key,
		client: client.client,
	}
}

func (this Key) ArbitraryCommand(command string, arguments ...string) []string {
	args := make([]string,len(arguments)+1)
	args[0] = this.key
	copy(args[1:],arguments)
	r := redis.SendStr(this.client.Rw,command,args...)
	checkForError(r.Err)
	if r.Elem != nil {
		return []string{r.Elem.String()}
	}
	return r.StringArray()
}

func (this Key) Exists() bool {
	exists, err := this.client.Exists(this.key)
	checkForError(err)
	return exists
}

func (this Key) Delete() bool {
	deleted, err := this.client.Del(this.key)
	checkForError(err)
	return deleted > 0
}

func (this Key) Type() string {
	t, err := this.client.Type(this.key)
	checkForError(err)
	return t
}

func (this Key) MoveTo(other Key) {
	err := this.client.Rename(this.key, other.key)
	checkForError(err)
}

func (this Key) MoveToIfEmpty(other Key) bool {
	moved, err := this.client.Renamenx(this.key, other.key)
	checkForError(err)
	return moved
}

func (this Key) ExpireIn(seconds int64) bool {
	set, err := this.client.Expire(this.key, seconds)
	checkForError(err)
	return set
}

func (this Key) ExpireAt(timestamp time.Time) bool {
	set, err := this.client.Expireat(this.key, timestamp.Unix())
	checkForError(err)
	return set
}

func (this Key) TimeToLive() int64 {
	ttl, err := this.client.Ttl(this.key)
	checkForError(err)
	return ttl
}
