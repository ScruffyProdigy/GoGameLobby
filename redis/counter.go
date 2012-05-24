package redis

type Counter string

func (this Counter) Set(val int64) {
	err := Client.Set(string(this), val)
	checkForError(err)
}

func (this Counter) SetIfEmpty(val int64) bool {
	set, err := Client.Setnx(string(this), val)
	checkForError(err)
	return set
}

func (this Counter) Get() int64 {
	val, err := Client.Get(string(this))
	checkForError(err)
	return val.Int64()
}

func (this Counter) GetSet(val int64) int64 {
	newVal, err := Client.Getset(string(this), val)
	checkForError(err)
	return newVal.Int64()
}

func (this Counter) Delete() bool {
	deleted, err := Client.Del(string(this))
	checkForError(err)
	return deleted > 0
}

func (this Counter) Clear() (int64, bool) {
	val := this.Get()
	return val, this.Delete()
}

func (this Counter) Increment() int64 {
	val, err := Client.Incr(string(this))
	checkForError(err)
	return val
}

func (this Counter) IncrementBy(val int64) int64 {
	newVal, err := Client.Incrby(string(this), val)
	checkForError(err)
	return newVal
}

func (this Counter) Decrement() int64 {
	val, err := Client.Decr(string(this))
	checkForError(err)
	return val
}

func (this Counter) DecrementBy(val int64) int64 {
	newVal, err := Client.Decrby(string(this), val)
	checkForError(err)
	return newVal
}
