package redis

type Set struct {
	Key
}

func newSet(client Redis, key string) Set {
	return Set{
		Key: newKey(client, key),
	}
}

func (this Set) IsValid() bool {
	thistype, err := this.client.Type(this.key)
	checkForError(err)
	return thistype == "set"
}

func (this Set) Add(item string) bool {
	isNew, err := this.client.Sadd(this.key, item)
	checkForError(err)
	return isNew
}

func (this Set) Remove(item string) bool {
	removed, err := this.client.Srem(this.key, item)
	checkForError(err)
	return removed
}

func (this Set) Members() []string {
	members, err := this.client.Smembers(this.key)
	checkForError(err)
	return members.StringArray()
}

func (this Set) IsMember(item string) bool {
	isMember, err := this.client.Sismember(this.key, item)
	checkForError(err)
	return isMember
}

func (this Set) Cardinality() int64 {
	cardinality, err := this.client.Scard(this.key)
	checkForError(err)
	return cardinality
}

func (this Set) RandomMember() string {
	item, err := this.client.Srandmember(this.key)
	checkForError(err)
	return item.String()
}

func (this Set) Pop() string {
	item, err := this.client.Spop(this.key)
	checkForError(err)
	return item.String()
}

func (this Set) Intersection(otherSet Set) []string {
	intersection, err := this.client.Sinter(this.key, otherSet.key)
	checkForError(err)
	return intersection.StringArray()
}

func (this Set) Union(otherSet Set) []string {
	union, err := this.client.Sunion(this.key, otherSet.key)
	checkForError(err)
	return union.StringArray()
}

func (this Set) Difference(otherSet Set) []string {
	difference, err := this.client.Sdiff(this.key, otherSet.key)
	checkForError(err)
	return difference.StringArray()
}

func (this Set) StoreIntersectionIn(newSet Set, otherSet Set) int64 {
	size, err := this.client.Sinterstore(this.key, newSet.key, otherSet.key)
	checkForError(err)
	return size
}

func (this Set) StoreUnionIn(newSet Set, otherSet Set) int64 {
	size, err := this.client.Sunionstore(this.key, newSet.key, otherSet.key)
	checkForError(err)
	return size
}

func (this Set) StoreDifferenceIn(newSet Set, otherSet Set) int64 {
	size, err := this.client.Sdiffstore(this.key, newSet.key, otherSet.key)
	checkForError(err)
	return size
}

func (this Set) MoveMemberTo(newSet Set, item string) bool {
	moved, err := this.client.Smove(this.key, newSet.key, item)
	checkForError(err)
	return moved
}
