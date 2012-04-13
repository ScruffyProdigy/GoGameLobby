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

//Middleware is a function that takes an HTTP request, a map of variables, and a function call to get the next Middleware
//it is expected to return an HTTP status code, a map of headers, and the resulting HTML
type Middleware func(req *http.Request, vars Vars, next NextFunc) (status int, header http.Header, message []byte)

//the Rack stores all of the Middleware
type Rack interface {
	Add(Middleware) //	Add will add a middleware into the list; it is order specific, so make sure you're calling it in the right order
	Go(Connection)  //	Go will start the rack with all of the middleware that have been added
}

type rack struct {
	middleware []Middleware
}

func (this *rack) Add(middle Middleware) {
	this.middleware = append(this.middleware, middle)
}

func (this *rack) Go(conn Connection) {
	conn.Go(func(w http.ResponseWriter, r *http.Request) {
		vars := make(Vars)
		index := -1
		var next NextFunc
		next = func() (int, http.Header, []byte) {
			index++
			if index >= len(this.middleware) {
				return BlankResponse().Results()
			}
			//			s,h,m := this.middleware[index](r, vars, next)
			//			fmt.Fprint(log.DebugLog(),"\n Debug - Passing Down: ",s,h,",and ",len(m)," bytes")
			//			return s,h,m
			return this.middleware[index](r, vars, next)
		}

		//this is the result of all of the middleware
		//so, write it out to the response writer
		status, headers, message := next()
		for k, _ := range headers {
			w.Header().Set(k, headers.Get(k))
		}
		w.WriteHeader(status)
		w.Write(message)

	})
}

//Up is the public interface to the default Rack
var Up Rack

//NewRack() is available in case you want a Rack separate from the default Rack 
func NewRack() Rack {
	return new(rack)
}

func init() {
	Up = NewRack()
}
