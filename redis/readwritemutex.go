package redis

import (
	"../mutex"
)

type ReadWriteMutex struct {
	readers int
	write   mutex.Mutex
	Read    mutex.Mutex
	Write   mutex.Mutex
}

type writeMutex struct {
	*ReadWriteMutex
}

func lockAllReads(rw *ReadWriteMutex, finalAction func()) func() {
	i := 0
	var lockNextRead func()
	lockNextRead = func() {
		if i < rw.readers {
			i++
			rw.Read.Force(lockNextRead)
		} else {
			finalAction()
		}
	}
	return lockNextRead
}

func (this writeMutex) Try(action func()) bool {
	return this.write.Try(lockAllReads(this.ReadWriteMutex, action))
}

func (this writeMutex) Force(action func()) {
	this.write.Force(lockAllReads(this.ReadWriteMutex, action))
}

func RWMutex(key string, readers int) (*ReadWriteMutex, error) {
	rw := new(ReadWriteMutex)
	rw.readers = readers

	var err error
	rw.write, err = Mutex(key + ":Write")
	if err != nil {
		return nil, err
	}

	rw.Read, err = Semaphore(key+":Read", readers)
	if err != nil {
		return nil, err
	}

	rw.Write = writeMutex{rw}
	return rw, nil
}
