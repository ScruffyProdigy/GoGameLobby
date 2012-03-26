package routes

import "net/http"
import "strings"
import "../log"
import "fmt"

var Root RouteBranch

func init() {
	Root = newRouteList()
}

func Implement(this RouteBranch) {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		//take the path, and split it up into a list of directions ending with '/'
		parsedRoute := strings.Split(req.URL.Path, "/")
		newParsedRoute := make([]string, 0, len(parsedRoute)+1)
		for _, section := range parsedRoute {
			if section != "" {
				l := len(newParsedRoute)
				newParsedRoute = newParsedRoute[0 : l+1]
				newParsedRoute[l] = section
			}
		}
		l := len(newParsedRoute)
		parsedRoute = newParsedRoute[0 : l+1]
		parsedRoute[l] = "/"

		for i, section := range parsedRoute {
			fmt.Fprint(log.DebugLog(), "\n", i, ":", section)
		}

		//follow the directions to dig down through subdirectories until we find the current one
		var currentRoute RouteBranch = this
		vars := make(VariableList)
		defer func() {
			//if there are any errors handling the request, render a 500 page
			if r := recover(); r != nil {
				log.Info("500")
				res.WriteHeader(500)
				fmt.Fprint(res, "<html><head><title>Error</title></head><body><h1>Error</h1><p>", r.(error).Error(), "</p></body></html>")
			}
		}()
		for _, section := range parsedRoute {
			fmt.Fprint(log.DebugLog(), "\nLooking for:"+section)
			found := false

			subroutes := make(chan Router)
			go currentRoute.GetSubroutes(subroutes)
			for subroute := range subroutes {
				switch subroute.Route(section, req, vars) {
				case route_elsewhere:
					log.Debug("\nRouting Elsewhere")
				case route_continue:
					found = true
					log.Debug("\nRouting Through")
					currentRoute = subroute.(RouteBranch)
					break
				case route_here:
					log.Debug("\nFound it - Right Here")
					subroute.(RouteTerminal).HandleRequest(res, req, vars)
					return
				}
			}
			log.Debug("\nNext Level!")
			//if we can't find what we're looking for, render a 404 page
			if !found {
				//404
				res.WriteHeader(404)
				fmt.Fprint(res, "<html><head><title>Not Found</title></head><body><h1>404 - Not Found</h1><p>Keep Looking!</p></body></html>")

				log.Info("\n\n404 - Not Found\n\n")
				return
			}
		}
	})
	http.ListenAndServe(":3000", nil)
}
