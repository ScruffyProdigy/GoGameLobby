package controller

import (
	"../rack"
	"../routes"
	"net/http"
)

func init() {
	root := routes.Get("/", func(res routes.Responder, req *http.Request, vars rack.Vars) {
		res.Render("test")
	})
	routes.Root.AddRoute(root)
}
