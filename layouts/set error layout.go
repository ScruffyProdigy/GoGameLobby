package layouts

import (
	"../rack"
	"net/http"
	"../templater"
	"strconv"
)

/*
	Statuser is a Middleware tool that is to be used with Encapsulator and Templater

	It checks the status of the response, and sets the layout to status number if it is found
	if it is not found, it sets it to a more general layout (if found) such as 40x, or 50x
*/
func Statuser(layoutstr, folder string) rack.Middleware {
	return func(r *http.Request, vars map[string]interface{}, next rack.NextFunc) (status int, header http.Header, message []byte) {
		status, header, message = next()

		layout := strconv.Itoa(status)
		if templater.Available(folder+layout) {
			vars[layoutstr] = layout
			return
		}

		layout = strconv.Itoa(status/100) + "0x"
		if templater.Available(folder+layout) {
			vars[layoutstr] = layout
			return
		}

		return
	}
}

/*
	SetErrorLayout is the default version of Statuser.  It works well with AddLayout
	It sets the status code into the "Layout" variable, if the layout is found in the "layouts" folder
*/
var SetErrorLayout rack.Middleware

func init() {
	SetErrorLayout = Statuser("Layout", "layouts/")
}
