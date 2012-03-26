package controller

import "../routing"
import "fmt"
import "../log"

func init() {
	root := routes.Get("/", func(res routes.Response, req routes.Request, vars routes.VariableList) {
		log.Info("Made it into here")
		fmt.Fprint(res, "<html><head><title>Test</title></head><body>Hello World</body></html>")
	})
	routes.Root.AddRoute(root)
}
