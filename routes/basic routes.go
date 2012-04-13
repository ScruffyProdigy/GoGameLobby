package routes

import (
	"../rack"
	"net/http"
)

type basicRoute struct {
	method  string
	name    string
	handler HandlerFunc
}

func (this *basicRoute) Route(section string, req *http.Request, vars rack.Vars) int {
	if section == this.name && req.Method == this.method {
		return route_here
	}
	return route_elsewhere
}

func (this *basicRoute) HandleRequest(res Responder, req *http.Request, vars rack.Vars) {
	this.handler(res, req, vars)
}

func newBasicRoute(method, name string, handler HandlerFunc) RouteTerminal {
	route := new(basicRoute)
	route.method = method
	route.name = name
	route.handler = handler
	return route
}

//Get provides a RouteTerminal that will direct a GET request to a specified handler
func Get(name string, handler HandlerFunc) RouteTerminal {
	return newBasicRoute("GET", name, handler)
}

//Post provides a RouteTermianl that will direct a POST request to a specified handler
func Post(name string, handler HandlerFunc) RouteTerminal {
	return newBasicRoute("POST", name, handler)
}

//Put provides a RouteTerminal that will direct a PUT request to a specified handler
func Put(name string, handler HandlerFunc) RouteTerminal {
	return newBasicRoute("PUT", name, handler)
}

//Delete provides a RouteTerminal that will direct a DELETE request to specified handler
func Delete(name string, handler HandlerFunc) RouteTerminal {
	return newBasicRoute("DELETE", name, handler)
}
