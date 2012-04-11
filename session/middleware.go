package session

import (
	"../rack"
	"net/http"
)

/*
	Middleware is the Middleware function that inserts a Session variable as "Session" into Rack variables
	This allows all later Middleware to have persistent effects
*/
func Middleware(r *http.Request, vars map[string]interface{}, next rack.NextFunc) (int, http.Header, []byte) {
	session := get(r)
	vars["session"] = Session(session)
	w := rack.CreateResponse(next())
	session.save(w)
	return w.Results()
}
