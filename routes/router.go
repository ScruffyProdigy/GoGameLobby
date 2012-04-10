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

func RouterWare(root RouteBranch) rack.Middleware {
	return func(r *http.Request, vars map[string]interface{}, next rack.NextFunc) (status int, header http.Header, message []byte) {

		parsedRoute := vars["parsedRoute"].([]string)
		currentRoute := root

		w := rack.BlankResponse()
		w2 := createResponder(w, vars)

		//in production we want this so that the user will get a 500 screen letting them know something went wrong
		//in development, we get more information if we just let it crash
		defer func() {
			rec := recover()
			if rec != nil {
				status = http.StatusInternalServerError
				message = []byte(getErrorString(rec))
			}
		}()

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
