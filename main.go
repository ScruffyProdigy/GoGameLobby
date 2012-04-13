package main

import (
	_ "./controllers"
	"./layouts"
	"./rack"
	"./routes"
	"./session"
	"./templater"
	//	"./errorhandler"
	"./facebooker"
	"./interceptor"
	"./notfound"
	"./oauther"
	"fmt"
)

func main() {
	i := interceptor.CreateInterceptor()

	facebooker.SetConfiguration("115772051792384", "211481baf989b0ac6ab4345debab6d91", "http://localhost:3000/", "facebook/start/", "facebook/authorize/", []string{})
	oauther.RegisterOauth(i, facebooker.Default)

	templater.LoadTemplates("./views")

	rack.Up.Add(layouts.AddLayout)
	rack.Up.Add(layouts.SetErrorLayout)
	rack.Up.Add(layouts.Defaulter("Layout", "base"))
	//	rack.Up.Add(errorhandler.ErrorHandler)	//in debug version, it's more useful to just let it crash, so we can get more error information
	rack.Up.Add(session.Middleware)
	rack.Up.Add(i.Middleware())
	rack.Up.Add(routes.Parser)
	rack.Up.Add(routes.RouterWare(routes.Root))
	rack.Up.Add(notfound.NotFound)

	fmt.Print("\n\nStarting!\n\n\n")

	rack.Up.Go(rack.HttpConnection(":3000"))
}
