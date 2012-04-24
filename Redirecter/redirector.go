package redirecter

import (
	"../log"
	"../rack"
	"net/http"
)

type Redirecter struct {
	Apply []rack.VarFunc
	Path  string
}

func (this Redirecter) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	return Go(r, vars, this.Path, this.Apply...)
}

func Go(r *http.Request, vars rack.Vars, path string, apply ...rack.VarFunc) (status int, header http.Header, message []byte) {
	log.Info("Redirecting to " + path)

	for _, a := range apply {
		vars.Apply(a)
	}

	w := rack.BlankResponse()
	http.Redirect(w, r, path, http.StatusFound)
	return w.Results()
}
