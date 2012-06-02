package redis

type Bits struct {
	Key
}

func newBits(client Root, key string) Bits {
	return Bits{
		Key: newKey(client, key),
	}
}

func (this Bits) IsValid() bool {
	thistype, err := this.client.Type(this.key)
	checkForError(err)
	return thistype == "string"
}

func (this Bits) SetTo(index int64, on bool) bool {
	var bit int
	if on {
		bit = 1
	}
	original, err := this.client.Setbit(this.key,index,bit)
	checkForError(err)
	return original == 1
}

func (this Bits) On(index int) bool {
	original, err := this.client.Setbit(this.key,index,1)
	checkForError(err)
	return original == 1
}

func (this Bits) Off(index int) bool {
	original, err := this.client.Setbit(this.key,index,0)
	checkForError(err)
	return original == 1
}

func (this Bits) Get(index int) bool {
	val, err := this.client.Getbit(this.key,index)
	checkForError(err)
	return val == 1
}

func (this Bits) Count(start, end int) int {
	s := this.ArbitraryCommand("BITCOUNT")
	return atoi(s[0])
}

func (this Bits) And(otherKey, resultKey Bits) int64 {
	r := this.client.SendStr(this.client.Rw,"BITOP","AND",resultKey.key,this.key,otherKey.key)
	return r.Int64()
}

func (this Bits) Or(otherKey, resultKey Bits) int64 {
	r := this.client.SendStr(this.client.Rw,"BITOP","OR",resultKey.key,this.key,otherKey.key)
	return r.Int64()
}

func (this Bits) Xor(otherKey, resultKey Bits) int64 {
	r := this.client.SendStr(this.client.Rw,"BITOP","XOR",resultKey.key,this.key,otherKey.key)
	return r.Int64()
}

func (this Bits) Not(resultKey Bits) int64 {
	r := this.client.SendStr(this.client.Rw,"BITOP","NOT",resultKey.key,this.key)
	return r.Int64()
}