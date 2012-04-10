package controller

import "../routes"
import "net/http"

func init() {
	root := routes.Get("/", func(res routes.Responder, req *http.Request, vars map[string]interface{}) {
		vars["Layout"] = "none"
		vars["Title"] = "Testing - 1 - 2 - 3"
		res.Render("test")
	})
	routes.Root.AddRoute(root)
}
