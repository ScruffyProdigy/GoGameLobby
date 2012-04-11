package routes

import "net/http"

type RoutingStatus int

const (
	route_error = iota
	route_elsewhere
	route_continue
	route_here
)

//a HandlerFunc is a function that multiple parts of the routing system use to dispatch control back to you once they've found the correct controller
type HandlerFunc func(Responder, *http.Request, map[string]interface{})

//a Router is a piece that helps to figure out where to send requests off to
//each Router is either a RouteTerminal or a RouteBranch
type Router interface {
	Route(section string, req *http.Request, vars map[string]interface{}) RoutingStatus // we send you the part of the request we're looking at, you tell us whether we're on the right path (and whether we can stop already)
}

//a Route Terminal is an end piece; there are no more routes to look for once we get here
type RouteTerminal interface {
	Router
	HandleRequest(res Responder, req *http.Request, vars map[string]interface{}) //you've told us this is the correct place to route a request, so, here is everything you need to respond to it
}

//a Route Branch is not an end piece, but will have more routes underneath it that it will direct us to
type RouteBranch interface {
	Router
	AddRoute(Router)                // AddRoute adds a subroute to look through once a request has been routed to this point
	GetSubroutes(out chan<- Router) //will sequentially send out all subroutes through the channel
}
