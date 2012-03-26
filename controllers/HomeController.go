package controller

import "../routing"
import "fmt"
import "../log"
import "net/http"

func init() {
	root := routes.Get("/", func(res http.ResponseWriter, req *http.Request, vars map[string]interface{}) {
		log.Info("Made it into here")
		fmt.Fprint(res, "<html><head><title>Test</title></head><body>Hello World</body></html>")
	})
	routes.Root.AddRoute(root)
}
