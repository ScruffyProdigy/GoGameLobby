package routes

import (
	"../rack"
	"net/http"
)

var Root *Router = NewRouter()

type Signaler interface {
	Run(r *http.Request, vars rack.Vars) bool
}

type Router struct {
	subroutes []*Router
	Action    rack.Middleware
	Routing   Signaler
}

func NewRouter() *Router {
	this := new(Router)
	this.subroutes = make([]*Router, 0)
	return this
}

func (this *Router) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	if vars.Apply(CurrentSection) == "/" {
		return this.Action.Run(r, vars, next)
	}
	for _, subroute := range this.subroutes {
		if subroute.Routing.Run(r, vars) {
			vars.Apply(nextSection)
			return subroute.Run(r, vars, next)
		}
	}
	return next()
}

func (this *Router) AddRoute(r ...*Router) {
	this.subroutes = append(this.subroutes, r...)
}
