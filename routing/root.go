package routes

import "net/http"
import "strings"
import "../log"

var Root RouteBranch

func init() {
	Root = newRouteList()
}

func Implement(this RouteBranch) {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		//take the path, and split it up into a list of directions ending with '/'
		parsedRoute := strings.Split(req.URL.Path, "/")
		var end int = len(parsedRoute)
		if parsedRoute[end-1] == "" {
			end -= 1
		}
		parsedRoute[end] = "/"
		if parsedRoute[0] == "" {
			parsedRoute = parsedRoute[1:]
		}

		//follow the directions to dig down through subdirectories until we find the current one
		var currentRoute RouteBranch = this
		vars := make(VariableList)
		defer func() {
			//if there are any errors handling the request, render a 500 page
			if r := recover(); r != nil {
				log.Info("500")
				//500
			}
		}()
		for _, section := range parsedRoute {

			found := false

			subroutes := make(chan Router)
			go currentRoute.GetSubroutes(subroutes)
			for subroute := range subroutes {
				switch subroute.Route(section, req, vars) {
				case route_elsewhere:
				case route_continue:
					currentRoute = subroute.(RouteBranch)
					break
				case route_here:
					subroute.(RouteTerminal).HandleRequest(res, req, vars)
					return
				}
			}
			//if we can't find what we're looking for, render a 404 page
			if !found {
				//404
				log.Info("\n\n404 - Not Found\n\n")
				return
			}
		}
	})
	http.ListenAndServe(":3000", nil)
}
