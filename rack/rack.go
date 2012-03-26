package rack

import "net/http"

type NextFunc func()

type Middleware func(http.ResponseWriter, *http.Request, map[string]interface{}, NextFunc)

type Interface interface {
	Add(Middleware)
	Go(Connection)
}

type implementation struct {
	middleware []Middleware
}

func (this *implementation) Add(middle Middleware) {
	if cap(this.middleware) == 0 {
		this.middleware = make([]Middleware, 0, 8)
	}

	n := len(this.middleware)
	for n+1 >= cap(this.middleware) {
		newmiddleware := make([]Middleware, n, 2*n)
		copy(newmiddleware, this.middleware)
		this.middleware = newmiddleware
	}

	this.middleware = this.middleware[0 : n+1]
	this.middleware[n] = middle
}

func (this *implementation) Go(conn Connection) {
	conn.Go(func(w http.ResponseWriter, r *http.Request) {
		vars := make(map[string]interface{})
		index := -1
		var next NextFunc
		next = func() {
			index++
			if index >= len(this.middleware) {
				return
			}
			this.middleware[index](w, r, vars, next)
		}
		next()
	})
}

var Up Interface

func init() {
	Up = new(implementation)
}
