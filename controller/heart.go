package controller

import (
	"../rack"
	"net/http"
	"../redirecter"
	"../templater"
	"../session"
)

// When Creating a Controller, you MUST put an anonymous controller.Heart into your controller (unless you really know what you're doing)
// Not only do some of the functions require some a couple of the default methods
type Heart struct {
	m ModelMap
	R *http.Request
	Vars rack.Vars
	Next rack.NextFunc
}

// this is how we hide the rack variables from the controllers who don't really care so much about these
// they are later accessible in case you need them, but for the most part, you can just ignore these
func (this *Heart) SetRackFuncVars(m ModelMap,r *http.Request, vars rack.Vars) {
	this.m = m
	this.R = r
	this.Vars = vars
}

// this sets the default response that the controller will give
// by default, this will lead to rendering for any GET requests, and redirection for PUT and POST
// DELETE is currently not implemented
func (this *Heart) SetDefaultResponse(next rack.NextFunc) {
	this.Next = next
}

// if you are calling another Rack Middleware, you should call this function to get the variables it will need to run
func (this Heart) GetRackFuncVars() (r *http.Request, vars rack.Vars, next rack.NextFunc) {
	return this.R,this.Vars,this.Next
}

// this should be the return value for most of your control function
// by default, this will lead to rendering for any GET requests, and redirection for PUT and POST
// DELETE is currently not implemented
func (this Heart) DefaultResponse() Response {
	return FromRack(this.Next())
}

// this should be the return value for most of your Create control functions
// you should pass it the resource you created
// since the default variable wasn't set because we didn't get into a specific resource
// this will set the default variable to the resource you just created
func (this Heart) RespondWith(object interface{}) Response {
	this.Vars[this.m.VarName()] = object
	return this.DefaultResponse()
}

// if things don't go according to plan, you can redirect somewhere else
// return this instead of DefaultResponse along with where you want to redirect to
func (this Heart) Redirection(url string) Response{
	return FromRack(redirecter.Go(this.R,this.Vars,url))
}

// use this if you want to render something other than the default template
// return this instead of DefaultResponse along with the template you want to render
func (this Heart) Rendering(tmpl string) Response{
	return FromRack(templater.Render(tmpl,this.Vars))
}

// this will get the form value from the form that was passed in
func (this Heart) GetFormValue(value string) string {
	return this.R.FormValue(value)
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

// if you want to apply any rack.VarFunc's, you can use this function
// look at the rack documentation if you don't know what this means
func (this Heart) Apply(f rack.VarFunc) interface{} {
	return this.Vars.Apply(f)
}

// this will get you a session variable
// look at the session documentation to figure out what you can do with it
func (this Heart) Session() session.Session {
	return this.Vars["session"].(session.Session)
}

// this will add a flash to the list of flashes
// flashes are typically used before redirecting
// all flashes stored before redirecting are available only after the next redirect, and are then immediately erased
func (this Heart) AddFlash(flash string) {
	this.Apply(session.AddFlash(flash))
}

// this will get access to all of the flashes stored before the last redirect
// they are also accessible within a template via {{.flashes}}
func (this Heart) GetFlashes() []string {
	return this.Vars["flashes"].([]string)
}