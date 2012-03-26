package routes

import "net/http"

type Response http.ResponseWriter
type Request *http.Request
type Variable interface{}
type VariableList map[string] Variable
type RoutingStatus int
const (
	route_error = iota
	route_elsewhere
	route_continue
	route_here
)
type HandlerFunc func(Response,Request,VariableList)


type Router interface{
	Route(section string, req Request, vars VariableList) RoutingStatus
}

type RouteTerminal interface{
	Router
	HandleRequest(res Response, req Request, vars VariableList)
}

type RouteBranch interface{
	Router
	AddRoute(Router);
	GetSubroutes(out chan<- Router)
}
