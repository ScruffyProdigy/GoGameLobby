package redis

type List struct {
	Key
}

func newList(client Redis, key string) List {
	return List{
		Key: newKey(client, key),
	}
}

func (this List) IsValid() bool {
	thistype, err := this.client.Type(this.key)
	checkForError(err)
	return thistype == "list"
}

func (this List) Length() int64 {
	length, err := this.client.Llen(this.key)
	checkForError(err)
	return length
}

func (this List) LeftPush(items ...string) (length int64) {
	var err error
	for _, item := range items {
		length, err = this.client.Lpush(this.key, item)
		checkForError(err)
	}
	return length
}

func (this List) LeftPushIfExists(item string) int64 {
	length, err := this.client.Lpushx(this.key, item)
	checkForError(err)
	return length
}

func (this List) RightPush(items ...string) (length int64) {
	var err error
	for _, item := range items {
		length, err = this.client.Rpush(this.key, item)
		checkForError(err)
	}
	return length
}

func (this List) RightPushIfExists(item string) int64 {
	length, err := this.client.Rpushx(this.key, item)
	checkForError(err)
	return length
}

func (this List) LeftPop() (string, bool) {
	item, err := this.client.Lpop(this.key)
	checkForError(err)
	return item.String(), item != nil
}

func (this List) BlockUntilLeftPop() string {
	item, _ := this.BlockUntilLeftPopWithTimeout(0)
	return item
}

func (this List) BlockUntilLeftPopWithTimeout(timeout int64) (string, bool) {
	item, err := this.client.Blpop([]string{this.key}, timeout)
	if isTimeout(err) {
		return "", false
	}
	checkForError(err)
	return item.StringMap()[this.key], true
}

func (this List) RightPop() (string, bool) {
	item, err := this.client.Rpop(this.key)
	checkForError(err)
	return item.String(), item != nil
}

func (this List) BlockUntilRightPop() string {
	item, _ := this.BlockUntilRightPopWithTimeout(0)
	return item
}

func (this List) BlockUntilRightPopWithTimeout(timeout int64) (string, bool) {
	item, err := this.client.Brpop([]string{this.key}, timeout)
	if isTimeout(err) {
		return "", false
	}
	checkForError(err)
	return item.StringMap()[this.key], true
}

func (this List) Index(index int) string {
	item, err := this.client.Lindex(this.key, index)
	checkForError(err)
	return item.String()
}

func (this List) Remove(item string) int64 {
	count, err := this.client.Lrem(this.key, 0, item)
	checkForError(err)
	return count
}

func (this List) Set(index int, item string) {
	err := this.client.Lset(this.key, index, item)
	checkForError(err)
}

func (this List) InsertBefore(pivot, item string) int64 {
	return this.Insert("BEFORE", pivot, item)
}

func (this List) InsertAfter(pivot, item string) int64 {
	return this.Insert("AFTER", pivot, item)
}

func (this List) Insert(location string, pivot, item string) int64 {
	index, err := this.client.Linsert(this.key, location, pivot, item)
	checkForError(err)
	return index
}

func (this List) GetFromRange(left, right int) []string {
	items, err := this.client.Lrange(this.key, left, right)
	checkForError(err)
	return items.StringArray()
}

func (this List) TrimToRange(left, right int) {
	err := this.client.Ltrim(this.key, left, right)
	checkForError(err)
}

func (this List) MoveLastItemToList(newList List) string {
	item, err := this.client.Rpoplpush(this.key, newList.key)
	checkForError(err)
	return item.String()
}

func (this List) BlockUntilMoveLastItemToList(newList List) string {
	item, _ := this.BlockUntilMoveLastItemToListWithTimeout(newList, 0)
	return item
}

func (this List) BlockUntilMoveLastItemToListWithTimeout(newList List, timeout int64) (string, bool) {
	item, err := this.client.Brpoplpush(this.key, newList.key, timeout)
	if isTimeout(err) {
		return "", false
	}
	checkForError(err)
	return item.String(), true
}
