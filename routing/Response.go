package routes

import (
	"../templater"
	"fmt"
	"net/http"
	"strings"
)

type Responder interface {
	RedirectTo(redirecttome Urler)
	Render(template string)
	RedirectWithCode(code int, redirecttome Urler)
	RenderWithCode(code int, template string)
}

type responder struct {
	w    http.ResponseWriter
	vars map[string]interface{}
}

type Urler interface {
	Url() []string
}

func GetUrl(this Urler) string {
	return "/" + strings.Join(this.Url(), "/")
}

func createResponder(w http.ResponseWriter, vars map[string]interface{}) *responder {
	r := new(responder)
	r.w = w
	r.vars = vars
	return r
}

func (this *responder) RedirectWithCode(code int, redirecttome Urler) {
	url := GetUrl(redirecttome)
	this.w.Header().Add("Location", url)
	this.w.WriteHeader(code)
	fmt.Fprint(this.w, "You are being redirected to ", url)
}

func (this *responder) RedirectTo(redirectome Urler) {
	this.RedirectWithCode(http.StatusFound, redirectome)
}

func (this *responder) Render(tmpl string) {
	this.RenderWithCode(http.StatusOK, tmpl)
}

func (this *responder) RenderWithCode(code int, tmpl string) {
	this.w.WriteHeader(code)
	var body bytes
	t := templater.Get(tmpl)
	t.Execute(&body, this.vars)
	this.vars["Body"] = (string)(body)
}

type bytes []byte

func (this *bytes) Write(p []byte) (n int, err error) {
	for _, c := range p {
		*this = append(*this, c)
		n++
	}
	return
}
