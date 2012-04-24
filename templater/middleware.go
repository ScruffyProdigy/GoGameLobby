package templater

import (
	"../rack"
	"net/http"
)

type Templater struct {
	Template string
}

func (this Templater) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	return Render(this.Template, vars)
}

func Render(s string, vars rack.Vars) (status int, header http.Header, message []byte) {
	w := rack.BlankResponse()
	Get(s).Execute(w, vars)
	return w.Results()
}
