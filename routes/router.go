/*
	Routes provides a basic router that operates on top of the Rack
	It Requires that the Parser Middleware is run before it
*/

package routes

import (
	"../log"
	"../rack"
	"net/http"
)

func getErrorString(rec interface{}) string {
	err, isError := rec.(error)
	str, isString := rec.(string)

	if isError {
		return err.Error()
	} else if isString {
		return str
	}
	return "Unknown Error"
}

/*
	RouterWare is the function that will create a basic router
	It requires that Parser is run before it
	Send RouterWare the root to your directory structure, and RouterWare will direct it to the correct Controller
*/
func RouterWare(root RouteBranch) rack.Middleware {
	return func(r *http.Request, vars map[string]interface{}, next rack.NextFunc) (status int, header http.Header, message []byte) {

		w := rack.BlankResponse()
		w2 := createResponder(w, vars)
		status, header, message = w.Results()

		//in production we want this so that the user will get a 500 screen letting them know something went wrong
		//in development, we get more information if we just let it crash
		defer func() {
			/*
						rec := recover()
			/*/
			var rec interface{} = nil
			/**/
			if rec != nil {
				status = http.StatusInternalServerError
				message = []byte(getErrorString(rec))
			}
		}()

		parsedRoute := vars["parsedRoute"].([]string)
		if parsedRoute == nil {
			log.Fatal("Please run the Parser middleware before running the RouterWare")
			panic("Please run the Parser middleware before running the RouterWare")
		}

		currentRoute := root

		//find the correct route, and have it handle the request
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
					subroute.(RouteTerminal).HandleRequest(w2, r, vars)
					next()
					return w.Results()
				}
			}
			//if we can't find what we're looking for, render a 404 page
			if !found {
				status = http.StatusNotFound
				return
			}
		}
		log.Unknown("False Exit in RouterWare")
		return
	}
}
