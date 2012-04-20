package controller

import (
	"../rack"
	"../routes"
	"../templater"
	"net/http"
)

func init() {
	routes.Root.Action = rack.Func(func(req *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
		w := rack.BlankResponse()
		t := templater.Get("test")
		t.Execute(w, vars)
		return w.Results()
	})
}
