package rack

import (
	"net/http"
)

type Vars map[string]interface{}
type VarFunc func(Vars) interface{}

func (this Vars) Apply(f VarFunc) interface{} {
	return f(this)
}

func NewVars() Vars {
	return make(Vars)
}

func (this Vars) Run(r *http.Request, vars Vars, next NextFunc) (status int, header http.Header, message []byte) {
	for k, v := range this {
		vars[k] = v
	}
	return next()
}
