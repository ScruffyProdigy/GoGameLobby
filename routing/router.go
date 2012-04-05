package routes

import "net/http"
import "../log"
import "fmt"
import "../rack"

func EndWare(root RouteBranch) rack.Middleware {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]interface{}, next rack.NextFunc) {
		defer func() {
			//if there are any errors handling the request, render a 500 page
			if rec := recover(); rec != nil {
				var error_string string
				log.Info("500")
				w.WriteHeader(500)

				fmt.Fprint(w, "<html><head><title>Error</title></head><body><h1>Error</h1>")

				err, isError := rec.(error)
				str, isString := rec.(string)

				if isError {
					error_string = err.Error()
				} else if isString {
					error_string = str
				} else {
					error_string = "Unknown Error"
				}

				fmt.Fprint(w, "<p>", error_string, "</p></body></html>")

			}
		}()
		var parsedRoute = vars["parsedRoute"].([]string)
		var currentRoute = root
		for _, section := range parsedRoute {
			found := false

			subroutes := make(chan Router)
			go currentRoute.GetSubroutes(subroutes)
			for subroute := range subroutes {
				switch subroute.Route(section, r, vars) {
				case route_elsewhere:
				case route_continue:
					found = true
					currentRoute = subroute.(RouteBranch)
					break
				case route_here:
					subroute.(RouteTerminal).HandleRequest(w, r, vars)
					return
				}
			}
			//if we can't find what we're looking for, render a 404 page
			if !found {
				//404
				w.WriteHeader(404)
				fmt.Fprint(w, "<html><head><title>Not Found</title></head><body><h1>404 - Not Found</h1><p>Keep Looking!</p></body></html>")
				return
			}
		}
	}
}
