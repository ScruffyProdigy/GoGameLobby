/*
	The Rack package allows you to break down your web server into smaller functions (called Middleware)

	Each Middleware gets inserted into a Rack, and then rack will proceed to call them one by one.
	Each one returns the combined results to the previous one, which can then be adjusted or altered

	This works as an alternative to simply writing to the provided http.ResponseWriter
	The advantage to using this method is an easy way to abstract away different parts of the program,
	which allows us to easily reuse smaller parts of the program in new websites.
	Also, we now have the ability to make adjustments to other parts of the program:
	Once something has been written to the ResponseWriter, there's no way to undo it, or to even change the headers.
	Middleware will frequently change the responses handed down by later Middleware, or write headers even after we know what most of the response will be
*/
package rack

import (
	"net/http"
)

//NextFunc is a function that gets passed to each of the Middleware so that they can interact with the next piece of Middleware
type NextFunc func() (int, http.Header, []byte)

//This is the signature of our interface.  Anything that fits this signature can work inside of a rack
type Func func(req *http.Request, vars Vars, next NextFunc) (status int, header http.Header, message []byte)

type Middleware interface {
	Run(req *http.Request, vars Vars, next NextFunc) (status int, header http.Header, message []byte)
}

func (this Func) Run(req *http.Request, vars Vars, next NextFunc) (status int, header http.Header, message []byte) {
	return this(req, vars, next)
}

type Rack []Middleware

func NewRack() *Rack {
	rack := make(Rack, 0, 2)
	return &rack
}

func (this *Rack) Add(m Middleware) {
	*this = append(*this, m)
}

func (this Rack) Run(r *http.Request, vars Vars, next NextFunc) (status int, header http.Header, message []byte) {
	index := -1
	var ourNext NextFunc
	ourNext = func() (int, http.Header, []byte) {
		index++
		if index >= len(this) {
			return next()
		}
		//			s,h,m := this.middleware[index](r, vars, next)
		//			fmt.Fprint(log.DebugLog(),"\n Debug - Passing Down: ",s,h,",and ",len(m)," bytes")
		//			return s,h,m
		return this[index].Run(r, vars, ourNext)
	}
	return ourNext()
}

//Up is the public interface to the default Rack
var Up *Rack = NewRack()

func NotFound() (status int, header http.Header, message []byte) {
	return http.StatusNotFound, make(http.Header), []byte("")
}

func Run(c Connection, m Middleware) error {
	return c.Go(func(w http.ResponseWriter, r *http.Request) {
		vars := NewVars()
		status, headers, message := m.Run(r, vars, NotFound)
		for k, _ := range headers {
			w.Header().Set(k, headers.Get(k))
		}
		w.WriteHeader(status)
		w.Write(message)
	})
}

func init() {
	Up = NewRack()
}
