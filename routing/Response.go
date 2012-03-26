package routes

import "net/http"
import "strings"
import "fmt"

type Urler interface {
	Url() []string
}

type Renderer interface {
	Data() []byte
}

func RedirectWithCode(this Response, code int, redirecttome Urler) {
	url := "/" + strings.Join(redirecttome.Url(), "/")
	this.Header().Add("Location", url)
	this.WriteHeader(code)
	fmt.Fprint(this, "You are being redirected to ", url)
}

func RenderWithCode(this Response, code int, renderme Renderer) {
	this.WriteHeader(code)
	this.Write(renderme.Data())
}

func RedirectTo(this Response, redirectome Urler) {
	RedirectWithCode(this, http.StatusFound, redirectome)
}

func Render(this Response, renderme Renderer) {
	RenderWithCode(this, http.StatusOK, renderme)
}

func NotFound(this Response, renderme Renderer) {
	RenderWithCode(this, http.StatusNotFound, renderme)
}

func Error(this Response, renderme Renderer) {
	RenderWithCode(this, http.StatusInternalServerError, renderme)
}
