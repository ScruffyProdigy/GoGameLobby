package session

import (
	"../rack"
	"net/http"
)

func Middleware(r *http.Request, vars map[string]interface{}, next rack.NextFunc) (int, http.Header, []byte) {
	session := get(r)
	vars["session"] = Session(session)
	w := rack.CreateResponse(next())
	session.save(w)
	return w.Results()
}
