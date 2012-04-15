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
	"./googleplusser"
	"./interceptor"
	"./loadconfiguration"
	"./login"
	"./notfound"
	"./oauther"
	"fmt"
)

func main() {
	i := interceptor.CreateInterceptor()

	var facebookData facebooker.Data
	err := configurations.Load("facebook", &facebookData)
	if err != nil {
		panic(err)
	}
	facebooker.SetConfiguration(facebookData)
	oauther.RegisterOauth(i, facebooker.Default, login.Logger)

	var googleData googleplusser.Data
	err = configurations.Load("google", &googleData)
	if err != nil {
		panic(err)
	}
	googleplusser.SetConfiguration(googleData)
	oauther.RegisterOauth(i, googleplusser.Default, login.Logger)

	login.RegisterLogout(i, "/logout/")

	templater.LoadTemplates("./views")

	rack.Up.Add(layouts.AddLayout)
	rack.Up.Add(layouts.SetErrorLayout)
	rack.Up.Add(layouts.Defaulter("Layout", "base"))
	//	rack.Up.Add(errorhandler.ErrorHandler)	//in debug version, it's more useful to just let it crash, so we can get more error information
	rack.Up.Add(session.Middleware)
	rack.Up.Add(login.Middleware)
	rack.Up.Add(i.Middleware())
	rack.Up.Add(routes.Parser)
	rack.Up.Add(routes.RouterWare(routes.Root))
	rack.Up.Add(notfound.NotFound)

	fmt.Print("\n\nStarting!\n\n\n")

	rack.Up.Go(rack.HttpConnection(":3000"))
}
