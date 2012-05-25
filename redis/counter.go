package redis

type Counter struct {
	Key
}

func newCounter(client Redis, key string) Counter {
	return Counter{
		Key: newKey(client, key),
	}
}

func (this Counter) Set(val int64) {
	err := this.client.Set(this.key, val)
	checkForError(err)
}

func (this Counter) SetIfEmpty(val int64) bool {
	set, err := this.client.Setnx(this.key, val)
	checkForError(err)
	return set
}

func (this Counter) Get() int64 {
	val, err := this.client.Get(this.key)
	checkForError(err)
	return val.Int64()
}

func (this Counter) GetSet(val int64) int64 {
	newVal, err := this.client.Getset(this.key, val)
	checkForError(err)
	return newVal.Int64()
}

func (this Counter) Clear() (int64, bool) {
	val := this.Get()
	return val, this.Delete()
}

func (this Counter) Increment() int64 {
	val, err := this.client.Incr(this.key)
	checkForError(err)
	return val
}

func (this Counter) IncrementBy(val int64) int64 {
	newVal, err := this.client.Incrby(this.key, val)
	checkForError(err)
	return newVal
}

func (this Counter) Decrement() int64 {
	val, err := this.client.Decr(this.key)
	checkForError(err)
	return val
}

func (this Counter) DecrementBy(val int64) int64 {
	newVal, err := this.client.Decrby(this.key, val)
	checkForError(err)
	return newVal
}
