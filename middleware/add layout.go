package middleware

import (
	"../log"
	"../rack"
	"../templater"
	"html/template"
	"net/http"
)

func AddLayout(r *http.Request, vars map[string]interface{}, next rack.NextFunc) (status int, header http.Header, message []byte) {
	status, header, body := next()
	vars["Body"] = template.HTML(body)
	w := rack.CreateResponse(status, header, []byte(""))

	layout, castable := vars["Layout"].(string)
	if !castable {
		log.Warning("Couldn't find Layout - Using \"base\"")
		layout = "base"
	}

	L := templater.Get("layouts/" + layout)
	if L == nil {
		panic("Layout Not Found")
	}

	L.Execute(w, vars)

	return w.Results()
}
