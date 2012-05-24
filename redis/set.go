package redis

type Set string

func (this Set) IsValid() bool {
	thistype, err := Client.Type(string(this))
	checkForError(err)
	return thistype == "set"
}

func (this Set) Delete() bool {
	deleted, err := Client.Del(string(this))
	checkForError(err)
	return deleted > 0
}

func (this Set) Add(item string) bool {
	isNew, err := Client.Sadd(string(this), item)
	checkForError(err)
	return isNew
}

func (this Set) Remove(item string) bool {
	removed, err := Client.Srem(string(this), item)
	checkForError(err)
	return removed
}

func (this Set) Members() []string {
	members, err := Client.Smembers(string(this))
	checkForError(err)
	return members.StringArray()
}

func (this Set) IsMember(item string) bool {
	isMember, err := Client.Sismember(string(this), item)
	checkForError(err)
	return isMember
}

func (this Set) Cardinality() int64 {
	cardinality, err := Client.Scard(string(this))
	checkForError(err)
	return cardinality
}

func (this Set) RandomMember() string {
	item, err := Client.Srandmember(string(this))
	checkForError(err)
	return item.String()
}

func (this Set) Pop() string {
	item, err := Client.Spop(string(this))
	checkForError(err)
	return item.String()
}

func (this Set) Intersection(otherSet Set) []string {
	intersection, err := Client.Sinter(string(this), string(otherSet))
	checkForError(err)
	return intersection.StringArray()
}

func (this Set) Union(otherSet Set) []string {
	union, err := Client.Sunion(string(this), string(otherSet))
	checkForError(err)
	return union.StringArray()
}

func (this Set) Difference(otherSet Set) []string {
	difference, err := Client.Sdiff(string(this), string(otherSet))
	checkForError(err)
	return difference.StringArray()
}

func (this Set) StoreIntersectionIn(newSet Set, otherSet Set) int64 {
	size, err := Client.Sinterstore(string(this), string(newSet), string(otherSet))
	checkForError(err)
	return size
}

func (this Set) StoreUnionIn(newSet Set, otherSet Set) int64 {
	size, err := Client.Sunionstore(string(this), string(newSet), string(otherSet))
	checkForError(err)
	return size
}

func (this Set) StoreDifferenceIn(newSet Set, otherSet Set) int64 {
	size, err := Client.Sdiffstore(string(this), string(newSet), string(otherSet))
	checkForError(err)
	return size
}

func (this Set) MoveTo(newSet Set, item string) bool {
	moved, err := Client.Smove(string(this), string(newSet), item)
	checkForError(err)
	return moved
}
