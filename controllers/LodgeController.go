package controller

import (
	"../models/lodge"
	"../routes"
	"net/http"
	"../rack"
	"../login"
	"../session"
	"../log"
	"fmt"
)

var L = lodge.L

func init() {
	rest := map[string]routes.HandlerFunc{
		"index": func(res routes.Responder, req *http.Request, vars rack.Vars) {
			var lodges []lodge.Lodge
			err := L.AllLodges(&lodges)
			if err != nil {
				panic(err)
			}

			vars["Lodges"] = lodges
			vars["Title"] = "Mason Lodges"

			res.Render("lodges/index")
		},
		"show": func(res routes.Responder, req *http.Request, vars rack.Vars) {
			l := vars["Lodge"].(*lodge.Lodge)
			vars["Title"] = l.Name
			
			res.Render("lodges/show")
		},
		"new": func(res routes.Responder, req *http.Request, vars rack.Vars) {
			vars["Title"] = "Create a Mason Lodge"

			res.Render("lodges/new")
		},
		"create": func(res routes.Responder, req *http.Request, vars rack.Vars) {
			err := req.ParseForm()
			if err != nil {
				panic(err)
			}

			var l lodge.Lodge
			defer func(){
				rec := recover()
				if rec != nil {
					vars.Apply(session.AddFlash("You fucked something up, please try again"))
					res.RedirectTo(routes.Url("/lodges/new"))
				}
			}()

			l.Name = req.FormValue("Lodge[Name]")
			l.Masons = []string{vars.Apply(login.CurrentUser()).(string)}

			err = L.AddLodge(&l)
			if err != nil {
				panic(err)
			}

			res.RedirectTo(l)
		},
	}

	lodgeResource := routes.Resource(L, rest)

	routes.Root.AddRoute(lodgeResource)
}