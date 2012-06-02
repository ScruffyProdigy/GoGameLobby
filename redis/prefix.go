package redis

type Prefix interface {
	Key(key string) Key
    String(key string) String
    Integer(key string) Integer
    Set(key string) Set
    List(key string) List
    Mutex(key string) Mutex
    Semaphore(key string, count int) Mutex
    ReadWriteMutex(key string, readers int) *ReadWriteMutex
    Channel(key string) Channel
	Prefix(key string) Prefix
}

type prefix struct {
	parent Prefix
	root    string
}

func (this *prefix) Key(key string) Key {
	return this.parent.Key(this.root + key)
}

func (this *prefix) String(key string) String {
	return this.parent.String(this.root + key)
}

func (this *prefix) Integer(key string) Integer {
	return this.parent.Integer(this.root + key)
}

func (this *prefix) Set(key string) Set {
	return this.parent.Set(this.root + key)
}

func (this *prefix) List(key string) List {
	return this.parent.List(this.root + key)
}

func (this *prefix) Mutex(key string) Mutex {
	return this.parent.Mutex(this.root + key)
}

func (this *prefix) Semaphore(key string, count int) Mutex {
	return this.parent.Semaphore(this.root + key,count)
}

func (this *prefix) ReadWriteMutex(key string, readers int) *ReadWriteMutex {
	return this.parent.ReadWriteMutex(this.root + key,readers)
}

func (this *prefix) Channel(key string) Channel {
	return this.parent.Channel(this.root + key)
}

func (this *prefix) Prefix(key string) Prefix {
	return newPrefix(this,key)
}

func newPrefix(parent Prefix, key string) Prefix {
	p := new(prefix)
	p.parent = parent
	p.root = key
	return p
}