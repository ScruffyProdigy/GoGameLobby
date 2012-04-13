package controller

import (
	"../rack"
	"../routes"
	"net/http"
)

func init() {
	root := routes.Get("/", func(res routes.Responder, req *http.Request, vars rack.Vars) {
		vars["Layout"] = "none"
		vars["Title"] = "Testing - 1 - 2 - 3"
		res.Render("test")
	})
	routes.Root.AddRoute(root)
}
