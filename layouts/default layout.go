package layouts

import (
	"../rack"
	"net/http"
)

/*
	Defaulter will set a variable to a specified default value
	we will use it to set a default layout in main.go
*/

func Defaulter(key string,value interface{}) rack.Middleware {
	return func(r *http.Request, vars map[string]interface{}, next rack.NextFunc) (status int, header http.Header, message []byte) {
		vars[key] = value
		return next()
	}
}
