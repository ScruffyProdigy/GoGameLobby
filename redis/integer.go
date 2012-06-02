package redis

type Integer struct {
	Key
}

func newInteger(client Root, key string) Integer {
	return Integer{
		Key: newKey(client, key),
	}
}

func (this Integer) IsValid() bool {
	thistype, err := this.client.Type(this.key)
	checkForError(err)
	return thistype == "string"
}

func (this Integer) Set(val int64) {
	err := this.client.Set(this.key, val)
	checkForError(err)
}

func (this Integer) SetIfEmpty(val int64) bool {
	set, err := this.client.Setnx(this.key, val)
	checkForError(err)
	return set
}

func (this Integer) Get() int64 {
	val, err := this.client.Get(this.key)
	checkForError(err)
	return val.Int64()
}

func (this Integer) GetSet(val int64) int64 {
	newVal, err := this.client.Getset(this.key, val)
	checkForError(err)
	return newVal.Int64()
}


//TODO: Make this happen atomically
func (this Integer) Clear() (int64, bool) {
	val := this.Get()
	return val, this.Delete()
}

func (this Integer) Increment() int64 {
	val, err := this.client.Incr(this.key)
	checkForError(err)
	return val
}

func (this Integer) IncrementBy(val int64) int64 {
	newVal, err := this.client.Incrby(this.key, val)
	checkForError(err)
	return newVal
}

func (this Integer) Decrement() int64 {
	val, err := this.client.Decr(this.key)
	checkForError(err)
	return val
}

func (this Integer) DecrementBy(val int64) int64 {
	newVal, err := this.client.Decrby(this.key, val)
	checkForError(err)
	return newVal
}
