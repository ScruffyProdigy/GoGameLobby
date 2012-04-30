package main

import (
	_ "./controllers"
	"./loadconfiguration"
	"./login"
	"./models"
	"fmt"
	"github.com/HairyMezican/Middleware/encapsulator"
	"github.com/HairyMezican/Middleware/errorhandler"
	"github.com/HairyMezican/Middleware/interceptor"
	"github.com/HairyMezican/Middleware/oauther"
	"github.com/HairyMezican/Middleware/oauther/facebooker"
	"github.com/HairyMezican/Middleware/oauther/googleplusser"
	"github.com/HairyMezican/Middleware/router"
	"github.com/HairyMezican/Middleware/sessioner"
	"github.com/HairyMezican/Middleware/statuser"
	"github.com/HairyMezican/TheRack/rack"
	"github.com/HairyMezican/TheTemplater/templater"
	"log"
	"os"
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
	cept := interceptor.New()

	//facebook
	fb := login.NewFacebooker(LoadFacebookData())
	oauther.SetIntercepts(cept, fb, login.HandleToken)

	//google plus
	gp := login.NewGooglePlusser(LoadGoogleData())
	oauther.SetIntercepts(cept, gp, login.HandleToken)

	//logging out
	cept.Intercept("/logout/", login.LogOut)

	//load the templates for the views
	templater.LoadFromFiles("./views", log.New(os.Stdout, "template - ", log.LstdFlags))

	//set up default variables
	defaults := rack.NewVars()
	defaults["Layout"] = "base"

	//set up the rack
	rack.Up.Add(defaults)
	rack.Up.Add(encapsulator.AddLayout)
	rack.Up.Add(statuser.SetErrorLayout)
	if mode != debug {
		rack.Up.Add(errorhandler.ErrorHandler) //in debug version, it's more useful to just let it crash, so we can get more error information
	}
	rack.Up.Add(sessioner.Middleware)
	rack.Up.Add(login.Middleware)
	rack.Up.Add(cept)
	rack.Up.Add(router.Parser)
	rack.Up.Add(router.Root)

	//alert the user as to where we are running
	var site = ":80"
	if mode == debug {
		site = ":3000"
	}

	fmt.Print("\n\nStarting at localhost" + site + "!\n\n\n")

	//We're ready to go!
	//run each request through the rack!
	err := rack.Run(rack.HttpConnection(site), rack.Up)
	if err != nil {
		fmt.Print("Error: " + err.Error())
	}
}
