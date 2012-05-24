package redis

type String string

func (this String) IsValid() bool {
	thistype, err := Client.Type(string(this))
	checkForError(err)
	return thistype == "string"
}

func (this String) Set(val string) {
	err := Client.Set(string(this), val)
	checkForError(err)
}

func (this String) SetIfEmpty(val string) bool {
	set, err := Client.Setnx(string(this), val)
	checkForError(err)
	return set
}

func (this String) Get() string {
	val, err := Client.Get(string(this))
	checkForError(err)
	return val.String()
}

func (this String) Delete() bool {
	deleted, err := Client.Del(string(this))
	checkForError(err)
	return deleted > 0
}

func (this String) Clear() (string, bool) {
	val := this.Get()
	return val, this.Delete()
}

func (this String) GetSet(val string) string {
	newVal, err := Client.Getset(string(this), val)
	checkForError(err)
	return newVal.String()
}

func (this String) Append(val string) int64 {
	newVal, err := Client.Append(string(this), val)
	checkForError(err)
	return newVal
}

func (this String) Length() int64 {
	length, err := Client.Strlen(string(this))
	checkForError(err)
	return length
}
