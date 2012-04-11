package routes

import (
	"../rack"
	"net/http"
	"strings"
)

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
