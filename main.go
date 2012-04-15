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
	"./login"
	"./notfound"
	"./oauther"
	"fmt"
)

func main() {
	i := interceptor.CreateInterceptor()

	facebooker.SetConfiguration(facebooker.Data{
		AppId:       "115772051792384",
		AppSecret:   "211481baf989b0ac6ab4345debab6d91",
		SiteUrl:     "http://localhost:3000/",
		AuthUrl:     "facebook/login/",
		RedirectUrl: "facebook/authorize/",
		Permissions: []string{}})
	oauther.RegisterOauth(i, facebooker.Default, login.Logger)

	googleplusser.SetConfiguration(googleplusser.Data{
		ClientID:     "588791846385.apps.googleusercontent.com",
		ClientSecret: "qrby3ReqJApabRfh-HBB1LWR",
		SiteUrl:      "http://localhost:3000/",
		StartUri:     "google/login/",
		RedirectUri:  "google/authorize",
		Permissions:  []string{googleplusser.UserPermission}})
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
