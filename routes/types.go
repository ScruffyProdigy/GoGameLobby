package routes

import "net/http"

type RoutingStatus int

const (
	route_error = iota
	route_elsewhere
	route_continue
	route_here
)

type HandlerFunc func(Responder, *http.Request, map[string]interface{})

type Router interface {
	Route(section string, req *http.Request, vars map[string]interface{}) RoutingStatus
}

type RouteTerminal interface {
	Router
	HandleRequest(res Responder, req *http.Request, vars map[string]interface{})
}

type RouteBranch interface {
	Router
	AddRoute(Router)
	GetSubroutes(out chan<- Router)
}
