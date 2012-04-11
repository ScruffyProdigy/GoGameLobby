package routes

import (
	"../rack"
	"net/http"
	"strings"
)
/*
	parser breaks down the request's URL path into a slice of strings
	later middleware will use it to direct control
*/
func Parser(r *http.Request, vars map[string]interface{}, next rack.NextFunc) (int, http.Header, []byte) {
	parsedRoute := strings.Split(r.URL.Path, "/")
	newParsedRoute := make([]string, 0, len(parsedRoute)+1)
	for _, section := range parsedRoute {
		if section != "" {
			l := len(newParsedRoute)
			newParsedRoute = newParsedRoute[0 : l+1]
			newParsedRoute[l] = section
		}
	}
	l := len(newParsedRoute)
	newParsedRoute = newParsedRoute[0 : l+1]
	newParsedRoute[l] = "/"

	vars["parsedRoute"] = newParsedRoute

	return next()
}
