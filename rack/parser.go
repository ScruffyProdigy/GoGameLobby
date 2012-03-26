package rack

import "net/http"
import "strings"

func Parser(w http.ResponseWriter, r *http.Request, vars map[string]interface{}, next NextFunc) {
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

	next()
}
