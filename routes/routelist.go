package routes

import (
	"../rack"
	"net/http"
)

type routeList struct {
	routes []Router
}

func (this *routeList) Route(section string, req *http.Request, vars rack.Vars) int {
	return route_continue
}

func (this *routeList) GetSubroutes(out chan<- Router) {
	for _, route := range this.routes {
		out <- route
	}
	close(out)
}

func (this *routeList) AddRoute(newroute Router) {
	n := len(this.routes)
	if n+1 > cap(this.routes) {
		newroutes := make([]Router, n, 2*n)
		copy(newroutes, this.routes)
		this.routes = newroutes
	}
	this.routes = this.routes[0 : n+1]
	this.routes[n] = newroute
}

func newRouteList() *routeList {
	this := new(routeList)
	this.routes = make([]Router, 0, 5)
	return this
}
