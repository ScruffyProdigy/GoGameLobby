package controller

import (
	"github.com/HairyMezican/Middleware/redirecter"
	"github.com/HairyMezican/Middleware/renderer"
	"github.com/HairyMezican/Middleware/sessioner"
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheRack/rack"
	"net/http"
	"strings"
)

// When Creating a Controller, you MUST put an anonymous controller.Heart into your controller (unless you really know what you're doing)
// Not only do some of the functions require some a couple of the default methods
type Heart struct {
	m      ModelMap
	Vars   map[string]interface{}
	Finish func()
}

// this is how we hide the rack variables from the controllers who don't really care so much about these
// they are later accessible in case you need them, but for the most part, you can just ignore these
func (this *Heart) SetRackFuncVars(m ModelMap, vars map[string]interface{}) {
	this.m = m
	this.Vars = vars
}

// this sets the default response that the controller will give
// by default, this will lead to rendering for any GET requests, and redirection for PUT and POST
// DELETE is currently not implemented
func (this *Heart) SetFinish(action func()) {
	this.Finish = action
}

// if you are calling another Rack Middleware, you should call this function to get the variables it will need to run
func (this Heart) GetRackFuncVars() (vars map[string]interface{}, next func()) {
	return this.Vars, this.Finish
}

// this should be the return value for most of your Create control functions
// you should pass it the resource you created
// since the default variable wasn't set because we didn't get into a specific resource
// this will set the default variable to the resource you just created
func (this Heart) RespondWith(object interface{}) {
	this.Vars[this.m.VarName()] = object
	this.Finish()
}

// if you have a piece of middleware that you want to respond with
// return this instead of Finish along with the middleware you want to run
func (this Heart) FinishWithMiddleware(m rack.Middleware) {
	m.Run(this.GetRackFuncVars())
}

// if things don't go according to plan, you can redirect somewhere else
// return this instead of Finish along with where you want to redirect to
func (this Heart) RedirectTo(url string) {
	(redirecter.V)(this.Vars).Redirect(url)
}

// use this if you want to render something other than the default template
// return this instead of Finish along with the template you want to render
func (this Heart) Render(tmpl string) {
	if !strings.Contains(tmpl, "/") {
		tmpl = this.m.RouteName() + "/" + tmpl
	}
	(renderer.V)(this.Vars).Render(tmpl)
}

func (this Heart) AddFlash(flash string) {
	(sessioner.V)(this.Vars).AddFlash(flash)
}

func (this Heart) Session() sessioner.V {
	return (sessioner.V)(this.Vars)
}

// this will get the form value from the form that was passed in
func (this Heart) GetFormValue(value string) string {
	return (httper.V)(this.Vars).GetRequest().FormValue(value)
}

func (this Heart) NotAuthorized() {
	(httper.V)(this.Vars).Status(http.StatusUnauthorized)
}

// this is used to set variables
// for the most part, this is used so that the template will have access to more variables when rendering
func (this Heart) Set(k string, v interface{}) {
	this.Vars[k] = v
}

// this is used to get previously set variables
// the most common variable to get is the one we stored for you for all member methods
func (this Heart) Get(k string) interface{} {
	return this.Vars[k]
}
