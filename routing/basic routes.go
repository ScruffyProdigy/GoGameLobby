package routes

import "fmt"
import "../log"
import "net/http"

type basicRoute struct {
	method  string
	name    string
	handler HandlerFunc
}

func (this *basicRoute) Route(section string, req *http.Request, vars map[string]interface{}) RoutingStatus {
	fmt.Fprint(log.DebugLog(), "\nComparing \"", section, "\" with \"", this.name, "\"")
	fmt.Fprint(log.DebugLog(), "\nComparing \"", req.Method, "\" with \"", this.method, "\"")
	if section == this.name && req.Method == this.method {
		return route_here
	}
	return route_elsewhere
}

func (this *basicRoute) HandleRequest(res http.ResponseWriter, req *http.Request, vars map[string]interface{}) {
	this.handler(res, req, vars)
}

func newBasicRoute(method, name string, handler HandlerFunc) RouteTerminal {
	route := new(basicRoute)
	route.method = method
	route.name = name
	route.handler = handler
	return route
}

func Get(name string, handler HandlerFunc) RouteTerminal {
	return newBasicRoute("GET", name, handler)
}

func Post(name string, handler HandlerFunc) RouteTerminal {
	return newBasicRoute("POST", name, handler)
}

func Put(name string, handler HandlerFunc) RouteTerminal {
	return newBasicRoute("PUT", name, handler)
}

func Delete(name string, handler HandlerFunc) RouteTerminal {
	return newBasicRoute("DELETE", name, handler)
}
