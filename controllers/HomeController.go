package controller

import "../routing"
import "../log"
import "net/http"

func init() {
	root := routes.Get("/", func(res routes.Responder, req *http.Request, vars map[string]interface{}) {
		log.Info("Made it into here")
		vars["layout"] = "none"
		vars["title"] = "Testing - 1 - 2 - 3"
		res.Render("test")
	})
	routes.Root.AddRoute(root)
}
