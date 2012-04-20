package controller

import (
	"../login"
	"../models/lodge"
	"../rack"
	"../redirecter"
	"../routes"
	"../session"
	"../templater"
	"net/http"
)

var L = lodge.L

func init() {
	rest := map[string]rack.Middleware{
		"index": rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
			w := rack.BlankResponse()

			var lodges []lodge.Lodge
			err := L.AllLodges(&lodges)
			if err != nil {
				panic(err)
			}

			vars["Lodges"] = lodges
			vars["Title"] = "Mason Lodges"

			templater.Get("lodges/index").Execute(w, vars)
			return w.Results()
		}),
		"show": rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
			w := rack.BlankResponse()

			l := vars["Lodge"].(*lodge.Lodge)
			vars["Title"] = l.Name

			templater.Get("lodges/show").Execute(w, vars)
			return w.Results()
		}),
		"new": rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
			w := rack.BlankResponse()

			vars["Title"] = "Create a Mason Lodge"

			templater.Get("lodges/new").Execute(w, vars)
			return w.Results()
		}),
		"create": rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
			w := rack.BlankResponse()

			err := r.ParseForm()
			if err != nil {
				panic(err)
			}

			var l lodge.Lodge
			defer func() {
				rec := recover()
				if rec != nil {
					reroute := redirecter.Go("/lodges/new",
						session.AddFlash("You fucked something up, please try again"))
					status, header, message = reroute.Run(r, vars, next)
				}
			}()

			l.Name = r.FormValue("Lodge[Name]")
			l.Masons = []string{vars.Apply(login.CurrentUser()).(string)}

			err = L.AddLodge(&l)
			if err != nil {
				panic(err)
			}

			http.Redirect(w, r, l.Url(), http.StatusFound)
			return w.Results()
		}),
	}

	lodgeResource := routes.Resource(L, rest)

	routes.Root.AddRoute(lodgeResource.Collection)
}
