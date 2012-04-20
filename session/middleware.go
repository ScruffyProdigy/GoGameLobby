package session

import (
	"../rack"
	"net/http"
)

/*
	Middleware is the Middleware function that inserts a Session variable as "Session" into Rack variables
	This allows all later Middleware to have persistent effects
*/
var Middleware = rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
	session := get(r)
	vars["session"] = Session(session)

	vars["flashes"] = session.Clear("flash")
	_, isStrings := vars["flashes"].([]string)
	if !isStrings {
		vars["flashes"] = []string{}
	}

	w := rack.CreateResponse(next())
	session.save(w)
	return w.Results()
})

func Set(k, v interface{}) rack.VarFunc {
	return func(vars rack.Vars) interface{} {
		vars["session"].(Session).Set(k, v)
		return nil
	}
}

func Get(k interface{}) rack.VarFunc {
	return func(vars rack.Vars) interface{} {
		return vars["session"].(Session).Get(k)
	}
}

func Clear(k interface{}) rack.VarFunc {
	return func(vars rack.Vars) interface{} {
		return vars["session"].(Session).Clear(k)
	}
}

func AddFlash(s string) rack.VarFunc {
	return func(vars rack.Vars) interface{} {
		a, isStrings := vars["session"].(Session).Get("flash").([]string)
		if !isStrings {
			a = []string{s}
		} else {
			a = append(a, s)
		}
		vars["session"].(Session).Set("flash", a)
		return nil
	}
}
