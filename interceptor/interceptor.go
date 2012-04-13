/*
	The interceptor package creates a Middleware that does a lightweight lookup for a bunch of static URLs
*/
package interceptor

import (
	"../log"
	"../rack"
	"net/http"
)

type Interceptor interface {
	Intercept(string, rack.Middleware)
	Middleware() rack.Middleware
}

type interceptor map[string]rack.Middleware

func (this interceptor) Intercept(url string, exec rack.Middleware) {
	if this[url] != nil {
		log.Error("Interception already registered!")
		panic("Interception already registered!")
	}
	this[url] = exec
}

func (this interceptor) Middleware() rack.Middleware {
	return func(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
		url := r.URL.Path
		exec := this[url]
		if exec == nil {
			return next()
		}
		return exec(r, vars, next)
	}
}

func CreateInterceptor() Interceptor {
	return make(interceptor)
}
