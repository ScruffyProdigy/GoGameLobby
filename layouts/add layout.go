package layouts

import (
	"../rack"
	"../templater"
	"html/template"
	"net/http"
)

/*
	Encapsulator is a Middleware to be used with templater
	it will encapsulate the current body, within a specified template
*/

func Encapsulator(layoutstr string, bodystr string, folder string) rack.Middleware {
	return func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
		status, header, body := next()

		layout, castable := vars[layoutstr].(string)
		if !castable {
			return
		}

		vars[bodystr] = template.HTML(body)
		w := rack.CreateResponse(status, header, []byte(""))

		L := templater.Get(folder + layout)
		if L == nil {
			panic("Layout Not Found")
		}

		L.Execute(w, vars)

		return w.Results()
	}
}

/*
	AddLayout is the default version of Encapsulator

	It will encapsulate the current body, within whichever template is in the "Layout" variable
	The layout will be found in the "layouts" folder, and will use {{.Body}} to specify the old body
*/
var AddLayout rack.Middleware

func init() {
	AddLayout = Encapsulator("Layout", "Body", "layouts/")
}
