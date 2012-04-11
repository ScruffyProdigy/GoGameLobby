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

/*
	RouterWare is the function that will create a basic router
	It requires that Parser is run before it
	Send RouterWare the root to your directory structure, and RouterWare will direct it to the correct Controller
*/
func RouterWare(root RouteBranch) rack.Middleware {
	return func(r *http.Request, vars map[string]interface{}, next rack.NextFunc) (status int, header http.Header, message []byte) {

		w := rack.BlankResponse()
		w2 := createResponder(w, vars)

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
					return w.Results()
				}
			}
			//if we can't find what we're looking for, perhaps the next Middleware will be able to
			if !found {
				return next()
			}
		}
		log.Unknown("False Exit in RouterWare")
		return
	}
}
