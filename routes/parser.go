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
var Parser = rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
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
	vars["currentSection"] = 0

	return next()
})

func CurrentSection(vars rack.Vars) interface{} {
	return vars["parsedRoute"].([]string)[vars["currentSection"].(int)]
}

func nextSection(vars rack.Vars) (result interface{}) {
	result = vars.Apply(CurrentSection)
	vars["currentSection"] = vars["currentSection"].(int) + 1
	return
}
