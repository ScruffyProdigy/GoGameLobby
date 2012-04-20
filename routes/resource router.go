package routes

import (
	"../models"
	"../rack"
	"net/http"
)

/*
a ResourceRouter assumes that it represents a RESTful resource, and will process it as such
it also allows you to add non-RESTful member and collection routes by exposing a route branch for each
*/
type ResourceRouter struct {
	Collection *Router //you can add non-RESTful collection-level routes here
	Member     *Router //you can add non-RESTful member-level routes here
}

type splitter struct {
	get, post, put, delete rack.Middleware
}

func (this splitter) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	var result rack.Middleware
	switch r.Method {
	case "GET":
		result = this.get
	case "POST":
		result = this.post
	case "PUT":
		result = this.put
	case "DELETE":
		result = this.delete
	default:
		log.Warning("Request with method:" + r.Method)
		return http.StatusBadRequest, make(http.Header), []byte("")
	}
	if result == nil {
		return next()
	}
	return result.Run(r, vars, next)
}

type memberSignaler struct {
	varName string
	indexer func(s string) interface{}
}

func (this memberSignaler) Run(r *http.Request, vars rack.Vars) bool {
	id := vars.Apply(currentSection).(string)
	result := this.indexer(id)
	if result == nil {
		return false
	}

	vars[this.varName] = result
	return true
}

type collectionSignaler struct {
	name string
}

func (this collectionSignaler) Run(r *http.Request, vars rack.Vars) bool {
	section := vars.Apply(currentSection).(string)
	if section == this.name {
		return true
	}
	return false
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
func Resource(m models.Collection, restfuncs map[string]rack.Middleware) *ResourceRouter {
	resource := new(ResourceRouter)

	resource.Member = NewRouter()
	resource.Member.routing = memberSignaler{varName: m.VarName(), indexer: func(s string) interface{} {
		return m.Indexer(s)
	}}
	resource.Member.Action = splitter{get: restfuncs["show"], put: restfuncs["update"], delete: restfuncs["destroy"]}
	if restfuncs["edit"] != nil {
		resource.Member.AddRoute(Get("edit", restfuncs["edit"]))
	}

	resource.Collection = NewRouter()
	resource.Collection.routing = collectionSignaler{name: m.RouteName()}
	resource.Collection.Action = splitter{get: restfuncs["index"], post: restfuncs["create"]}
	if restfuncs["new"] != nil {
		resource.Collection.AddRoute(Get("new", restfuncs["new"]))
	}
	resource.Collection.AddRoute(resource.Member)

	return resource
}
