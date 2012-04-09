package routes

import (
	"../log"
	"../rack"
	"../templater"
	"fmt"
	"html/template"
	"net/http"
)

func RouteWare(root RouteBranch) rack.Middleware {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]interface{}, next rack.NextFunc) {
		//		defer handleErrors(w)	//in production we want this so that the user will get a 500 screen letting them know something went wrong
		//in development, we get more information if we just let it crash

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
	layout, castable := vars["Layout"].(string)
	if !castable {
		log.Warning("\nWarning: Couldn't find Layout - Using \"base\"")
		layout = "base"
	}

	_, castable = vars["Body"].(string)
	if !castable {
		log.Warning("\nWarning:Couldn't find Body - Using \"\"")
		vars["Body"] = ""
	}
	vars["Body"] = template.HTML(vars["Body"].(string))

	L := templater.Get("layouts/" + layout)
	if L == nil {
		log.Error("\nError: Layout Not Found - Printing body as is")
		//*///	** Start of Code Switch - One slash before the start uses second set of code, two slashes uses first set of code
		fmt.Fprint(w, vars["Body"].(string))
		/*/// ** Middle of Code Switch	
				panic("Layout Not Found")
		/**/ //	** End of Code Switch
	} else {
		L.Execute(w, vars)
	}
}
