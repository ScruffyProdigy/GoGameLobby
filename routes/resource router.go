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

/*
a ResourceRouter assumes that it represents a RESTful resource, and will process it as such
it also allows you to add non-RESTful member and collection routes by exposing a route branch for each
*/
type ResourceRouter struct {
	name            string
	collectionfuncs *routeList
	memberfuncs     *routeList

	//public
	Collection RouteBranch //you can add non-RESTful collection-level routes here
	Member     RouteBranch //you can add non-RESTful member-level routes here
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

/*
	Resource will return a RESTful Resource Router
	it expects
	name: a string that represents the name of the resource.  This is used in the routing process
	restfuncs: the RESTful routes that this resource expects to handle
		the usable keys in the map are: "index","new","create","show","edit","update", and "delete"
	variablename: If we are drilling down into a member of the resource, we will add a variable to the rack variables, and this will be the name that it will set
	getter:	if we need to get a member resource, you'll have to help us;  we'll give you the string representing the ID, you give us the resource
*/
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
