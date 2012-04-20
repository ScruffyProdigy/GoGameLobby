package routes

import (
	"../rack"
	"net/http"
)

type basicRoute struct {
	method string
	name   string
}

func (this *basicRoute) Run(req *http.Request, vars rack.Vars) bool {
	sec := vars.Apply(currentSection).(string)
	if sec == this.name && req.Method == this.method {
		return true
	}
	return false
}

//Get provides a RouteTerminal that will direct a GET request to a specified handler
func Get(name string, m rack.Middleware) (result *Router) {
	result = NewRouter()
	result.routing = &basicRoute{method: "GET", name: name}
	result.Action = m
	return
}

//Post provides a RouteTermianl that will direct a POST request to a specified handler
func Post(name string, m rack.Middleware) (result *Router) {
	result = NewRouter()
	result.routing = &basicRoute{method: "POST", name: name}
	result.Action = m
	return
}

//Put provides a RouteTerminal that will direct a PUT request to a specified handler
func Put(name string, m rack.Middleware) (result *Router) {
	result = NewRouter()
	result.routing = &basicRoute{method: "PUT", name: name}
	result.Action = m
	return
}

//Delete provides a RouteTerminal that will direct a DELETE request to specified handler
func Delete(name string, m rack.Middleware) (result *Router) {
	result = NewRouter()
	result.routing = &basicRoute{method: "DELETE", name: name}
	result.Action = m
	return
}
