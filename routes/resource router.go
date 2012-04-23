package routes

import (
	"../rack"
	"net/http"
	"../log"
	"../templater"
)

/*
a ResourceRouter assumes that it represents a RESTful resource, and will process it as such
it also allows you to add non-RESTful member and collection routes by exposing a route branch for each
*/
type ModelMap interface {
	Indexer(string) (interface{},bool)
	RouteName() string
	VarName() string
}

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
	indexer func(s string) (interface{},bool)
}

func (this memberSignaler) Run(r *http.Request, vars rack.Vars) bool {
	id := vars.Apply(currentSection).(string)
	result,found := this.indexer(id)
	if !found {
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


func Resource(m ModelMap) *ResourceRouter {
	resource := new(ResourceRouter)
	
	restfuncs := GetRestMap(m)

	resource.Member = NewRouter()
	resource.Member.routing = memberSignaler{varName: m.VarName(), indexer: func(s string) (interface{},bool) {
		return m.Indexer(s)
	}}
	memberRouter := splitter{}
	if restfuncs["show"] != nil {
		memberRouter.get = rack.Rack{restfuncs["show"],templater.Templater{m.RouteName()+"/show"}}
	}
	if restfuncs["update"] != nil {
		memberRouter.put = rack.Rack{restfuncs["update"]}
	}
	if restfuncs["destroy"] != nil {
		memberRouter.delete = rack.Rack{restfuncs["destroy"]}
	}
	resource.Member.Action = memberRouter
	
	if restfuncs["edit"] != nil {
		resource.Member.AddRoute(Get("edit", rack.Rack{restfuncs["edit"],templater.Templater{m.RouteName()+"/edit"}}))
	}

	resource.Collection = NewRouter()
	resource.Collection.routing = collectionSignaler{name: m.RouteName()}
	collectionRouter := splitter{}
	if restfuncs["index"] != nil {
		collectionRouter.get = rack.Rack{restfuncs["index"],templater.Templater{m.RouteName()+"/index"}} 		
	}
	if restfuncs["create"] != nil {
		collectionRouter.post = rack.Rack{restfuncs["create"]}		
	}
	resource.Collection.Action = collectionRouter
	
	if restfuncs["new"] != nil {
		resource.Collection.AddRoute(Get("new", rack.Rack{restfuncs["new"],templater.Templater{m.RouteName()+"/new"}}))
	}
	resource.Collection.AddRoute(resource.Member)

	return resource
}

func (this ResourceRouter) AddTo(superroute *Router) {
	superroute.AddRoute(this.Collection)
}