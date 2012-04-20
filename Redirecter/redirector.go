package redirecter

import (
	"../rack"
	"net/http"
)

type Redirecter struct {
	Apply []rack.VarFunc
	Path  string
}

func (this Redirecter) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	for _, a := range this.Apply {
		vars.Apply(a)
	}

	w := rack.BlankResponse()
	http.Redirect(w, r, this.Path, http.StatusFound)
	return w.Results()
}

func Go(path string, apply ...rack.VarFunc) Redirecter {
	return Redirecter{Path: path, Apply: apply}
}
