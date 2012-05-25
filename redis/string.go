package redis

type String struct {
	Key
}

func newString(client Redis, key string) String {
	return String{
		Key: newKey(client, key),
	}
}

func (this String) IsValid() bool {
	thistype, err := this.client.Type(this.key)
	checkForError(err)
	return thistype == "string"
}

func (this String) Set(val string) {
	err := this.client.Set(this.key, val)
	checkForError(err)
}

func (this String) SetIfEmpty(val string) bool {
	set, err := this.client.Setnx(this.key, val)
	checkForError(err)
	return set
}

func (this String) Get() string {
	val, err := this.client.Get(this.key)
	checkForError(err)
	return val.String()
}

func (this String) Clear() (string, bool) {
	val := this.Get()
	return val, this.Delete()
}

func (this String) Replace(val string) string {
	newVal, err := this.client.Getset(this.key, val)
	checkForError(err)
	return newVal.String()
}

func (this String) Append(val string) int64 {
	newVal, err := this.client.Append(this.key, val)
	checkForError(err)
	return newVal
}

func (this String) Length() int64 {
	length, err := this.client.Strlen(this.key)
	checkForError(err)
	return length
}
