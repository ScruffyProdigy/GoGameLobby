/*
	Not Found Package

	This middleware is intended to be used at the end of a Middleware chain
	It assumes that none of the other middleware found what they were looking for, and returns a not found error
*/
package notfound

import (
	"../rack"
	"net/http"
)

func NotFound(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	return http.StatusNotFound, make(http.Header), []byte("")
}
