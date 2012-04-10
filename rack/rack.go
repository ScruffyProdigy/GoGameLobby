package rack

import (
	"net/http"
)

type NextFunc func() (int, http.Header, []byte)

type Middleware func(*http.Request, map[string]interface{}, NextFunc) (int, http.Header, []byte)

type Interface interface {
	Add(Middleware)
	Go(Connection)
}

type implementation struct {
	middleware []Middleware
}

func (this *implementation) Add(middle Middleware) {
	this.middleware = append(this.middleware, middle)
}

func (this *implementation) Go(conn Connection) {
	conn.Go(func(w http.ResponseWriter, r *http.Request) {
		vars := make(map[string]interface{})
		index := -1
		var next NextFunc
		next = func() (int, http.Header, []byte) {
			index++
			if index >= len(this.middleware) {
				return BlankResponse().Results()
			}
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

var Up Interface

func init() {
	Up = new(implementation)
}
