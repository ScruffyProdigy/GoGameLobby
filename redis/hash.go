package redis

type Hash struct {
	Key
}

func newHash(client Root,key string) Hash {
	return Hash{
		Key: newKey(client, key),
	}
}

func (this Hash) IsValid() bool {
	thistype, err := this.client.Type(this.key)
	checkForError(err)
	return thistype == "hash"
}

func (this Hash) String(key string) HashString {
	return newHashString(this,key)
}

func (this Hash) Integer(key string) HashInteger {
	return newHashInteger(this,key)
}

func (this Hash) Size() int64 {
	size,err := this.client.Hlen(this.key)
	checkForError(err)
	return size
}

func (this Hash) Get() map[string]string {
	vals,err := this.client.Hgetall(this.key)
	checkForError(err)
	return vals.StringMap()
}

type HashField struct {
	parent Hash
	key string
}

func newHashField(hash Hash, key string) HashField {
	return HashField{
		parent: hash,
		key: key,
	}
}

func (this HashField) Delete() bool {
	deleted, err := this.parent.client.Hdel(this.parent.key,this.key)
	checkForError(err)
	return deleted
}

func (this HashField) Exists() bool {
	exists, err := this.parent.client.Hexists(this.parent.key,this.key)
	checkForError(err)
	return exists
}

type HashString struct {
	HashField
}

func newHashString(hash Hash, key string) HashString {
	return HashString{
		HashField:newHashField(hash,key),
	}
}

func (this HashString) Get() string {
	val, err := this.parent.client.Hget(this.parent.key,this.key)
	checkForError(err)
	return val.String()
}

func (this HashString) Set(val string) bool {
	isNew, err := this.parent.client.Hset(this.parent.key,this.key,val)
	checkForError(err)
	return isNew
}

func (this HashString) SetIfEmpty(val string) bool {
	isSet, err := this.parent.client.Hsetnx(this.parent.key,this.key,val)
	checkForError(err)
	return isSet
}

type HashInteger struct {
	HashField
}

func newHashInteger(hash Hash, key string) HashInteger {
	return HashInteger{
		HashField:newHashField(hash,key),
	}
}

func (this HashInteger) Get() int64 {
	val, err := this.parent.client.Hget(this.parent.key,this.key)
	checkForError(err)
	return val.Int64()
}

func (this HashInteger) Set(val int64) bool {
	isNew, err := this.parent.client.Hset(this.parent.key,this.key,val)
	checkForError(err)
	return isNew
}

func (this HashInteger) SetIfEmpty(val int64) bool {
	isSet, err := this.parent.client.Hsetnx(this.parent.key,this.key,val)
	checkForError(err)
	return isSet
}

func (this HashInteger) Increment(val int64) int64 {
	result, err := this.parent.client.Hincrby(this.parent.key,this.key,val)
	checkForError(err)
	return result
}

func (this HashInteger) Decrement(val int64) int64 {
	result, err := this.parent.client.Hincrby(this.parent.key,this.key,-val)
	checkForError(err)
	return result
}