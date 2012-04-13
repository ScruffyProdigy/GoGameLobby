package routes

import (
	"../rack"
	"../templater"
	"net/http"
)

/*
	Typically the controller would rather not worry about the presentation details
	The Responder encapsulates this, and just provides the interface that controllers will want to deal with
*/
type Responder interface {
	Render(template string)                        // Render will render a template to the user using a status of http.StatusOK
	RedirectTo(redirecttome Urler)                 //RedirectTo will send a redirect message back to the user using a status of http.StatusFound
	RenderWithCode(code int, template string)      // RenderWithCode is like Render, except you can use any HTTPStatus code
	RedirectWithCode(code int, redirecttome Urler) //RedirectWithCode is like RedirectTo, except you can use any HTTPStatus code
}

type responder struct {
	w    http.ResponseWriter
	vars rack.Vars
}

/*
	When we redirect, we need to know where to redirect to; 
	an Urler is something that has a Url() method that we can use to get the URL we need to redirect to
*/
type Urler interface {
	Url() string
}

type Url string

func (this Url) Url() string {
	return string(this)
}

func createResponder(w http.ResponseWriter, vars rack.Vars) *responder {
	r := new(responder)
	r.w = w
	r.vars = vars
	return r
}

func (this *responder) RedirectWithCode(code int, redirecttome Urler) {
	url := redirecttome.Url()

	this.w.Header().Add("Location", url)
	this.w.WriteHeader(code)
}

func (this *responder) RedirectTo(redirectome Urler) {
	this.RedirectWithCode(http.StatusFound, redirectome)
}

func (this *responder) Render(tmpl string) {
	this.RenderWithCode(http.StatusOK, tmpl)
}

func (this *responder) RenderWithCode(code int, tmpl string) {
	this.w.WriteHeader(code)
	t := templater.Get(tmpl)
	t.Execute(this.w, this.vars)
}
