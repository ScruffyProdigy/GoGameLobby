package controller

import (
	"github.com/HairyMezican/Middleware/router"
	"github.com/HairyMezican/TheRack/rack"
	"net/http"
)

/*
a ResourceRouter assumes that it represents a RESTful resource, and will process it as such
it also allows you to add non-RESTful member and collection routes by exposing a route branch for each
*/
type ModelMap interface {
	//These are provided by you
	RouteName() string                  // We need a way to know whether or not a url is trying to access your resource; tell us what route you want to lay claim to
	Indexer(string) (interface{}, bool) //We need some way to go from a collection to a specific resource provided by a url - we give you a string, and you tell us if the resource is there, and, if so, what the resource is
	VarName() string                    // We need somewhere to store the resource you give us, so you can access it later
	//These are provided by us
	SetRackFuncVars(ModelMap, *http.Request, rack.Vars)
	SetDefaultResponse(rack.Next)
}

type ControllerShell struct {
	Collection *router.Router //you can add non-RESTful collection-level routes here
	Member     *router.Router //you can add non-RESTful member-level routes here
}

type splitter struct {
	get, post, put, delete rack.Middleware
}

func (this splitter) Run(r *http.Request, vars rack.Vars, next rack.Next) (status int, header http.Header, message []byte) {
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
		return http.StatusBadRequest, make(http.Header), []byte("")
	}
	if result == nil {
		//that particular method wasn't set, but perhaps a later middleware will take care of it
		return next()
	}
	return result.Run(r, vars, next)
}

type memberSignaler struct {
	varName string
	indexer func(string) (interface{}, bool)
}

func (this memberSignaler) Run(r *http.Request, vars rack.Vars) bool {
	id := vars.Apply(router.CurrentSection).(string)
	result, found := this.indexer(id)
	if !found {
		return false
	}

	vars[this.varName] = result

	return true
}

type collectionSignaler struct {
	m    ModelMap
	name string
}

func (this collectionSignaler) Run(r *http.Request, vars rack.Vars) bool {
	this.m.SetRackFuncVars(this.m, r, vars)
	section := vars.Apply(router.CurrentSection).(string)
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

func AddMapRoutes(superroute *router.Router, routemap map[string]rack.Middleware, methodfinder func(string, rack.Middleware) *router.Router) {
	for name, action := range routemap {
		superroute.AddRoute(methodfinder(name, action))
	}
}

func AddMapListRoutes(superroute *router.Router, maplist mapList) {
	AddMapRoutes(superroute, maplist.get, router.Get)
	AddMapRoutes(superroute, maplist.put, router.Put)
	AddMapRoutes(superroute, maplist.post, router.Post)
	AddMapRoutes(superroute, maplist.delete, router.Delete)
	AddMapRoutes(superroute, maplist.all, router.All)
}

func RegisterController(m ModelMap) *ControllerShell {
	resource := new(ControllerShell)

	restfuncs := GetRestMap(m)
	memberfuncs := GetGenericMapList(m, "Member")
	collectionfuncs := GetGenericMapList(m, "Collection")

	resource.Member = router.NewRouter()
	resource.Member.Routing = memberSignaler{varName: m.VarName(), indexer: func(s string) (interface{}, bool) {
		return m.Indexer(s)
	}}
	memberRouter := splitter{}
	memberRouter.get = restfuncs["Show"]
	memberRouter.put = restfuncs["Update"]
	memberRouter.delete = restfuncs["Destroy"]
	resource.Member.Action = memberRouter

	if restfuncs["Edit"] != nil {
		memberfuncs.get["Edit"] = restfuncs["Edit"]
	}
	AddMapListRoutes(resource.Member, memberfuncs)

	resource.Collection = router.NewRouter()
	resource.Collection.Routing = collectionSignaler{m: m, name: m.RouteName()}
	collectionRouter := splitter{}
	collectionRouter.get = restfuncs["Index"]
	collectionRouter.post = restfuncs["Create"]
	resource.Collection.Action = collectionRouter

	if restfuncs["New"] != nil {
		collectionfuncs.get["New"] = restfuncs["New"]
	}
	AddMapListRoutes(resource.Collection, collectionfuncs)

	resource.Collection.AddRoute(resource.Member)

	return resource
}

func (this ControllerShell) AddTo(superroute *router.Router) {
	superroute.AddRoute(this.Collection)
}

func (this ControllerShell) AddToRoot() {
	router.Root.AddRoute(this.Collection)
}

func (this ControllerShell) AddAsSubresource(parent *ControllerShell) {
	parent.Member.AddRoute(this.Collection)
}

func (this ControllerShell) AddAsSubmethod(parent *ControllerShell) {
	parent.Collection.AddRoute(this.Collection)
}
