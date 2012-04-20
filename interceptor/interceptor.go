/*
	The interceptor package creates a Middleware that does a lightweight lookup for a bunch of static URLs
*/
package interceptor

import (
	"../log"
	"../rack"
	"net/http"
)

type Interceptor map[string]rack.Middleware

func (this Interceptor) Intercept(url string, exec rack.Middleware) {
	if this[url] != nil {
		log.Error("Interception '" + url + "' already registered!")
		panic("Interception already registered!")
	}
	this[url] = exec
}

func (this Interceptor) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
	url := r.URL.Path
	exec := this[url]
	if exec == nil {
		return next()
	}
	return exec.Run(r, vars, next)
}

func NewInterceptor() Interceptor {
	return make(Interceptor)
}
