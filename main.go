package main

import (
	_ "./controllers"
	"./errorhandler"
	"./facebooker"
	"./googleplusser"
	"./interceptor"
	"./layouts"
	"./loadconfiguration"
	"./log"
	"./login"
	"./models"
	"./oauther"
	"./rack"
	"./routes"
	"./session"
	"./templater"
	"fmt"
)

const (
	debug = iota
	release
)

const mode = debug

func LoadFacebookData() (result facebooker.Data) {
	err := configurations.Load("facebook", &result)
	if err != nil {
		panic(err)
	}
	return
}

func LoadGoogleData() (result googleplusser.Data) {
	err := configurations.Load("google", &result)
	if err != nil {
		panic(err)
	}
	return
}

func main() {

	//set up the models
	model.SetUp() //can't happen during models's init, because it needs to wait until each of the models has initialized

	//set up the interceptor routes
	cept := interceptor.NewInterceptor()

	//facebook
	facebooker.SetConfiguration(LoadFacebookData())
	oauther.RegisterOauth(cept, facebooker.Default, login.CreateHandler(facebooker.Default))

	//google plus
	googleplusser.SetConfiguration(LoadGoogleData())
	oauther.RegisterOauth(cept, googleplusser.Default, login.CreateHandler(googleplusser.Default))

	//logging out
	cept.Intercept("/logout/", login.LogOut)

	//load the templates for the views
	templater.LoadTemplates("./views")

	//set up default variables
	defaults := rack.NewVars()
	defaults["Layout"] = "base"

	//set up the rack
	rack.Up.Add(defaults)
	rack.Up.Add(layouts.AddLayout)
	rack.Up.Add(layouts.SetErrorLayout)
	if mode != debug {
		rack.Up.Add(errorhandler.ErrorHandler) //in debug version, it's more useful to just let it crash, so we can get more error information
	}
	rack.Up.Add(session.Middleware)
	rack.Up.Add(login.Middleware)
	rack.Up.Add(cept)
	rack.Up.Add(routes.Parser)
	rack.Up.Add(routes.Root)

	//alert the user as to where we are running
	var site = ":80"
	if mode == debug {
		site = ":3000"
	}

	fmt.Print("\n\nStarting at localhost" + site + "!\n\n\n")

	//set an appropriate logging level
	if mode == release {
		log.SetLogLevel(log.Level_Warning)
	}

	//We're ready to go!
	//run each request through the rack!
	err := rack.Run(rack.HttpConnection(site), rack.Up)
	if err != nil {
		fmt.Print("Error: " + err.Error())
	}
}
