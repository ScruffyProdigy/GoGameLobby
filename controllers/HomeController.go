package controller

import "../routing"
import "net/http"

func init() {
	root := routes.Get("/", func(res routes.Responder, req *http.Request, vars map[string]interface{}) {
		vars["layout"] = "none"
		vars["title"] = "Testing - 1 - 2 - 3"
		res.Render("test")
	})
	routes.Root.AddRoute(root)
}
