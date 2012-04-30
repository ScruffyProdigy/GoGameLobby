package controllers

import (
	"github.com/HairyMezican/Middleware/router"
	"github.com/HairyMezican/TheRack/rack"
	"github.com/HairyMezican/TheTemplater/templater"
	"net/http"
)

func init() {
	router.Root.Action = rack.Func(func(req *http.Request, vars rack.Vars, next rack.Next) (status int, header http.Header, message []byte) {
		w := rack.BlankResponse()
		t, _ := templater.Get("test")
		t.Execute(w, vars)
		return w.Results()
	})
}
