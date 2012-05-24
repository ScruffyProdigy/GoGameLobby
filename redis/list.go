package redis

type List string

func (this List) IsValid() bool {
	thistype, err := Client.Type(string(this))
	checkForError(err)
	return thistype == "list"
}

func (this List) Delete() bool {
	deleted, err := Client.Del(string(this))
	checkForError(err)
	return deleted > 0
}

func (this List) Length() int64 {
	length, err := Client.Llen(string(this))
	checkForError(err)
	return length
}

func (this List) LeftPush(items ...string) (length int64) {
	var err error
	for _, item := range items {
		length, err = Client.Lpush(string(this), item)
		checkForError(err)
	}
	return length
}

func (this List) LeftPushIfExists(item string) int64 {
	length, err := Client.Lpushx(string(this), item)
	checkForError(err)
	return length
}

func (this List) RightPush(items ...string) (length int64) {
	var err error
	for _, item := range items {
		length, err = Client.Rpush(string(this), item)
		checkForError(err)
	}
	return length
}

func (this List) RightPushIfExists(item string) int64 {
	length, err := Client.Rpushx(string(this), item)
	checkForError(err)
	return length
}

func (this List) LeftPop() string {
	item, err := Client.Lpop(string(this))
	checkForError(err)
	return item.String()
}

func (this List) BlockUntilLeftPop() string {
	item, _ := this.BlockUntilLeftPopWithTimeout(0)
	return item
}

func (this List) BlockUntilLeftPopWithTimeout(timeout int64) (string, error) {
	item, err := Client.Blpop([]string{string(this)}, timeout)
	if isTimeout(&err) {
		return "", err
	}
	checkForError(err)
	return item.StringMap()[string(this)], nil
}

func (this List) RightPop() string {
	item, err := Client.Rpop(string(this))
	checkForError(err)
	return item.String()
}

func (this List) BlockUntilRightPop() string {
	item, _ := this.BlockUntilRightPopWithTimeout(0)
	return item
}

func (this List) BlockUntilRightPopWithTimeout(timeout int64) (string, error) {
	item, err := Client.Brpop([]string{string(this)}, timeout)
	if isTimeout(&err) {
		return "", err
	}
	checkForError(err)
	return item.StringMap()[string(this)], nil
}

func (this List) Index(index int) string {
	item, err := Client.Lindex(string(this), index)
	checkForError(err)
	return item.String()
}

func (this List) Remove(item string) int64 {
	count, err := Client.Lrem(string(this), 0, item)
	checkForError(err)
	return count
}

func (this List) Set(index int, item string) {
	err := Client.Lset(string(this), index, item)
	checkForError(err)
}

func (this List) InsertBefore(pivot, item string) int64 {
	return this.Insert("BEFORE", pivot, item)
}

func (this List) InsertAfter(pivot, item string) int64 {
	return this.Insert("AFTER", pivot, item)
}

func (this List) Insert(location string, pivot, item string) int64 {
	index, err := Client.Linsert(string(this), location, pivot, item)
	checkForError(err)
	return index
}

func (this List) GetFromRange(left, right int) []string {
	items, err := Client.Lrange(string(this), left, right)
	checkForError(err)
	return items.StringArray()
}

func (this List) TrimToRange(left, right int) {
	err := Client.Ltrim(string(this), left, right)
	checkForError(err)
}

func (this List) MoveLastItemToList(newList List) string {
	item, err := Client.Rpoplpush(string(this), string(newList))
	checkForError(err)
	return item.String()
}

func (this List) BlockUntilMoveLastItemToList(newList List) string {
	item, _ := this.BlockUntilMoveLastItemToListWithTimeout(newList, 0)
	return item
}

func (this List) BlockUntilMoveLastItemToListWithTimeout(newList List, timeout int64) (string, error) {
	item, err := Client.Brpoplpush(string(this), string(newList), timeout)
	if isTimeout(&err) {
		return "", err
	}
	checkForError(err)
	return item.String(), nil
}
