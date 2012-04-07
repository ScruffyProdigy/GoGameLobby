package routes

import (
	"../log"
	"../rack"
	"../templater"
	"fmt"
	"net/http"
)

func RouteWare(root RouteBranch) rack.Middleware {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]interface{}, next rack.NextFunc) {
		//		defer handleErrors(w)

		parsedRoute := vars["parsedRoute"].([]string)
		currentRoute := root

		w2 := createResponder(w, vars)

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

func RenderWare(w http.ResponseWriter, r *http.Request, vars map[string]interface{}, next rack.NextFunc) {
	layout, castable := vars["layout"].(string)
	if !castable {
		log.Debug("\nCouldn't find Layout")
		layout = "base"
	}

	_, castable = vars["body"].(string)
	if !castable {
		log.Debug("\nCouldn't find Body")
		vars["body"] = ""
	}

	log.Debug("\nLayout: " + layout)
	log.Debug("\nTesting!")
	L := templater.Get("layouts/" + layout)
	fmt.Fprint(log.DebugLog(), "\nResult:", L)
	if L == nil {
		log.Debug("\nNot Found - printing body separate")
		fmt.Fprint(w, vars["body"].(string))
	} else {
		log.Debug("\nFound - rendering template")
		L.Execute(w, vars)
	}
}
