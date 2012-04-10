package middleware

import (
	"../rack"
	"net/http"
)

func SetErrorLayout(r *http.Request, vars map[string]interface{}, next rack.NextFunc) (status int, header http.Header, message []byte) {
	status, header, message = next()
	if status == 404 {
		vars["Layout"] = "404"
	} else if status/100 == 4 {
		vars["Layout"] = "400"
	} else if status/100 == 5 {
		vars["Layout"] = "500"
	}
	return
}
