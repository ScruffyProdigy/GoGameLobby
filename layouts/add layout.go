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
type Encapsulator struct {
	LayoutVar string //the variable this will look for to find the layout
	BodyVar   string //the variable this will set the old body into that the layout should look for to reapply
	Folder    string //the folder to look for the layouts in
}

func (this Encapsulator) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	status, header, body := next()

	layout, castable := vars[this.LayoutVar].(string)
	if !castable {
		//no "layout", just let it through
		return
	}

	vars[this.BodyVar] = template.HTML(body)
	w := rack.CreateResponse(status, header, []byte(""))

	L := templater.Get(this.Folder + "/" + layout)
	if L == nil {
		panic("Layout Not Found")
	}

	L.Execute(w, vars)

	return w.Results()
}

/*
	AddLayout is the default version of Encapsulator

	It will encapsulate the current body, within whichever template is in the "Layout" variable
	The layout will be found in the "layouts" folder, and will use {{.Body}} to specify the old body
*/
var AddLayout = Encapsulator{LayoutVar: "Layout", BodyVar: "Body", Folder: "layouts"}
