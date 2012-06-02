package redis

import (
	"strconv"
)

type SortedSet struct {
	Key
}

func newSortedSet(client Root, key string) Set {
	return Set{
		Key: newKey(client, key),
	}
}

func (this SortedSet) IsValid() bool {
	thistype, err := this.client.Type(this.key)
	checkForError(err)
	return thistype == "zset"
}

func (this SortedSet) Add(item string, score float64) bool {
	isNew, err := this.client.Zadd(this.key, score, item)
	checkForError(err)
	return isNew
}

func (this SortedSet) IncrementBy(item string, score float64) float64 {
	isNew, err := this.client.Zincrby(this.key, score, item)
	checkForError(err)
	return isNew
}

func (this SortedSet) Remove(item string) bool {
	removed, err := this.client.Zrem(this.key, item)
	checkForError(err)
	return removed
}

func (this SortedSet) Size() int64 {
	size, err := this.client.Zcard(this.key)
	checkForError(err)
	return size
}

func (this SortedSet) Count(min, max float64) int64 {
	count,err := this.client.Zcount(this.key,min,max)
	checkForError(err)
	return count
}

func (this SortedSet) IndexOf(item string) (int64,bool) {
	s := this.ArbitraryCommand("ZRANK",item)
	if len(s) == 0 {
		return 0,false
	}
	return atoi(s[0]),true
}

func (this SortedSet) ReverseIndexOf(item string) (int64,bool) {
	s := this.ArbitraryCommand("ZREVRANK",item)
	if len(s) == 0 {
		return 0,false
	}
	return atoi(s[0]),true
}

func (this SortedSet) ScoreOf(item string) (float64,bool) {
	s := this.ArbitraryCommand("ZSCORE",item)
	if len(s) == 0 {
		return 0,false
	}
	return atof(s[0]),true
}

func (this SortedSet) IndexedBetween(start, stop int) []string {
	vals,err := this.client.Zrange(this.key,start,stop)
	checkForError(err)
	return vals.StringArray()
}

func (this SortedSet) ReverseIndexedBetween(start, stop int) []string {
	vals,err := this.client.Zrevrange(this.key,start,stop)
	checkForError(err)
	return vals.StringArray()
}

func (this SortedSet) RemoveByIndex(start, stop int) int64 {
	count,err := this.client.Zremrangebyrank(this.key,start,stop)
	checkForError(err)
	return count
}

func (this SortedSet) RemoveScoresBetween(min, max float64) int64 {
	count,err := this.client.Zremrangebyscore(this.key,min,max)
	checkForError(err)
	return count
}

func (this SortedSet) RemoveScoresAbove(min float64) int64 {
	e := this.ArbitraryCommand("ZREMRANGEBYSCORE",ftoa(min),"+inf")
	return atoi(e[0])
}

func (this SortedSet) RemoveScoresBelow(max float64) int64 {
	e := this.ArbitraryCommand("ZREMRANGEBYSCORE",this.key,"-inf")
	return atoi(e[0])
}

func (this SortedSet) ScoredBetween(min, max float64) []string {
	vals,err := this.client.Zrangebyscore(this.key,ftoa(min),ftoa(max))
	checkForError(err)
	return vals.StringArray()
}

func (this SortedSet) ScoredAbove(min float64) []string {
	vals,err := this.client.Zrangebyscore(this.key,ftoa(min),"+inf")
	checkForError(err)
	return vals.StringArray()
}

func (this SortedSet) ScoredBelow(max float64) []string {
	vals,err := this.client.Zrangebyscore(this.key,"-inf",ftoa(max))
	checkForError(err)
	return vals.StringArray()
}

func (this SortedSet) ScoredBetweenWithLimit(min, max float64, offset, count int) []string {
	vals,err := this.client.Zrangebyscore(this.key,ftoa(min),ftoa(max), "LIMIT", itoa(offset), itoa(count))
	checkForError(err)
	return vals.StringArray()
}

func (this SortedSet) ScoredAboveWithLimit(min float64, offset, count int) []string {
	vals,err := this.client.Zrangebyscore(this.key,ftoa(min),"+inf", "LIMIT", itoa(offset), itoa(count))
	checkForError(err)
	return vals.StringArray()
}

func (this SortedSet) ScoredBelowWithLimit(max float64, offset, count int) []string {
	vals,err := this.client.Zrangebyscore(this.key,"-inf",ftoa(max), "LIMIT", itoa(offset), itoa(count))
	checkForError(err)
	return vals.StringArray()
}

func (this SortedSet) LowestScores(offset,count int) []string {
	vals,err := this.client.Zrangebyscore(this.key,"-inf","+inf","LIMIT",itoa(offset),itoa(count))
	checkForError(err)
	return vals.StringArray()
}

func (this SortedSet) DescendingScoredBetween(min, max float64) []string {
	return this.ArbitraryCommand("ZREVRANGEBYSCORE",ftoa(min),ftoa(max))
}

func (this SortedSet) DescendingScoredAbove(min float64) []string {
	return this.ArbitraryCommand("ZREVRANGEBYSCORE",ftoa(min),"+inf")
}

func (this SortedSet) DescendingScoredBelow(max float64) []string {
	return this.ArbitraryCommand("ZREVRANGEBYSCORE","-inf",ftoa(max))
}

func (this SortedSet) DescendingScoredBetweenWithLimit(min, max float64, offset, count int) []string {
	return this.ArbitraryCommand("ZREVRANGEBYSCORE",ftoa(min),ftoa(max), "LIMIT", itoa(offset), itoa(count))
}

func (this SortedSet) DescendingScoredAboveWithLimit(min float64, offset, count int) []string {
	return this.ArbitraryCommand("ZREVRANGEBYSCORE",ftoa(min),"+inf", "LIMIT", itoa(offset), itoa(count))
}

func (this SortedSet) DescendingScoredBelowWithLimit(max float64, offset, count int) []string {
	return this.ArbitraryCommand("ZREVRANGEBYSCORE","-inf",ftoa(max), "LIMIT", itoa(offset), itoa(count))
}

func (this SortedSet) HighestScores(offset, count int) []string {
	return this.ArbitraryCommand("-inf","+inf","LIMIT",itoa(offset),itoa(count))
}

func (this SortedSet) StoreIntersectionIn(newSet SortedSet, otherSet SortedSet) int64 {
	size, err := this.client.Zinterstore(this.key, []string{newSet.key}, otherSet.key)
	checkForError(err)
	return size
}

func (this SortedSet) StoreWeightedIntersectionIn(newSet SortedSet, otherSet SortedSet, weightA, weightB float64) int64 {
	size,err := this.client.Zinterstore(this.key, []string{newSet.key}, otherSet.key, "WEIGHTS", ftoa(weightA), ftoa(weightB))
	checkForError(err)
	return size
}

func (this SortedSet) StoreMinIntersectionIn(newSet SortedSet, otherSet SortedSet) int64 {
	size, err := this.client.Zinterstore(this.key, []string{newSet.key}, otherSet.key, "AGGREGATE", "MIN")
	checkForError(err)
	return size
}

func (this SortedSet) StoreWeightedMinIntersectionIn(newSet SortedSet, otherSet SortedSet, weightA, weightB float64) int64 {
	size,err := this.client.Zinterstore(this.key, []string{newSet.key}, otherSet.key, "WEIGHTS", ftoa(weightA), ftoa(weightB), "AGGREGATE", "MIN")
	checkForError(err)
	return size
}

func (this SortedSet) StoreMaxIntersectionIn(newSet SortedSet, otherSet SortedSet) int64 {
	size, err := this.client.Zinterstore(this.key, []string{newSet.key}, otherSet.key, "AGGREGATE", "MAX")
	checkForError(err)
	return size
}

func (this SortedSet) StoreWeightedMaxIntersectionIn(newSet SortedSet, otherSet SortedSet, weightA, weightB float64) int64 {
	size,err := this.client.Zinterstore(this.key, []string{newSet.key}, otherSet.key, "WEIGHTS", ftoa(weightA), ftoa(weightB), "AGGREGATE", "MAX")
	checkForError(err)
	return size
}

func (this SortedSet) StoreUnionIn(newSet SortedSet, otherSet SortedSet) int64 {
	size, err := this.client.Zunionstore(this.key, []string{newSet.key}, otherSet.key)
	checkForError(err)
	return size
}

func (this SortedSet) StoreWeightedUnionIn(newSet SortedSet, otherSet SortedSet, weightA, weightB float64) int64 {
	size,err := this.client.Zunionstore(this.key, []string{newSet.key}, otherSet.key, "WEIGHTS", ftoa(weightA), ftoa(weightB))
	checkForError(err)
	return size
}

func (this SortedSet) StoreMinUnionIn(newSet SortedSet, otherSet SortedSet) int64 {
	size, err := this.client.Zunionstore(this.key, []string{newSet.key}, otherSet.key, "AGGREGATE", "MIN")
	checkForError(err)
	return size
}

func (this SortedSet) StoreWeightedMinUnionIn(newSet SortedSet, otherSet SortedSet, weightA, weightB float64) int64 {
	size,err := this.client.Zunionstore(this.key, []string{newSet.key}, otherSet.key, "WEIGHTS", ftoa(weightA), ftoa(weightB), "AGGREGATE", "MIN")
	checkForError(err)
	return size
}

func (this SortedSet) StoreMaxUnionIn(newSet SortedSet, otherSet SortedSet) int64 {
	size, err := this.client.Zunionstore(this.key, []string{newSet.key}, otherSet.key, "AGGREGATE", "MAX")
	checkForError(err)
	return size
}

func (this SortedSet) StoreWeightedMaxUnionIn(newSet SortedSet, otherSet SortedSet, weightA, weightB float64) int64 {
	size,err := this.client.Zunionstore(this.key, []string{newSet.key}, otherSet.key, "WEIGHTS", ftoa(weightA), ftoa(weightB), "AGGREGATE", "MAX")
	checkForError(err)
	return size
}