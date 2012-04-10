package routes

import "net/http"

type memberRouter struct {
	variablename string
	getter       func(string) interface{}
	resource     RouteBranch
}

func (this *memberRouter) Route(section string, req *http.Request, vars map[string]interface{}) RoutingStatus {
	received := this.getter(section)
	if received == nil {
		return route_elsewhere
	}
	vars[this.variablename] = received
	return this.resource.Route(section, req, vars)
}

func (this *memberRouter) AddRoute(router Router) {
	this.resource.AddRoute(router)
}

func (this *memberRouter) GetSubroutes(out chan<- Router) {
	this.resource.GetSubroutes(out)
}

func newMemberRouter(variablename string, getter func(string) interface{}, resource RouteBranch) *memberRouter {
	this := new(memberRouter)
	this.variablename = variablename
	this.getter = getter
	this.resource = resource
	return this
}

type ResourceRouter struct {
	name            string
	collectionfuncs *routeList
	memberfuncs     *routeList

	//public
	Collection RouteBranch
	Member     RouteBranch
}

func (this *ResourceRouter) Route(section string, req *http.Request, vars map[string]interface{}) RoutingStatus {
	if this.name == section {
		return this.collectionfuncs.Route(section, req, vars)
	}
	return route_elsewhere
}

func (this *ResourceRouter) GetSubroutes(out chan<- Router) {
	this.Collection.GetSubroutes(out)
}

func (this *ResourceRouter) AddRoute(router Router) {
	this.Collection.AddRoute(router)
}

func Resource(name string, restfuncs map[string]HandlerFunc, variablename string, getter func(index string) interface{}) *ResourceRouter {
	resource := new(ResourceRouter)
	resource.name = name
	resource.collectionfuncs = newRouteList()
	resource.memberfuncs = newRouteList()
	resource.Collection = resource.collectionfuncs
	resource.Member = resource.memberfuncs

	resource.Collection.AddRoute(newMemberRouter(variablename, getter, resource.Member))

	var function HandlerFunc

	function = restfuncs["index"]
	if function != nil {
		resource.Collection.AddRoute(Get("/", function))
	}

	function = restfuncs["new"]
	if function != nil {
		resource.Collection.AddRoute(Get("new", function))
	}

	function = restfuncs["create"]
	if function != nil {
		resource.Collection.AddRoute(Post("/", function))
	}

	function = restfuncs["show"]
	if function != nil {
		resource.Member.AddRoute(Get("/", function))
	}

	function = restfuncs["edit"]
	if function != nil {
		resource.Member.AddRoute(Get("edit", function))
	}

	function = restfuncs["update"]
	if function != nil {
		resource.Member.AddRoute(Put("/", function))
	}

	function = restfuncs["delete"]
	if function != nil {
		resource.Member.AddRoute(Delete("/", function))
	}

	return resource
}
