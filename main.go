package main

import (
	"./clashgetter"
	"./controllers"
	"./defaulter"
	"./loadconfiguration"
	"./login"
	"./models"
	"./requestlogger"
	"./websocketcontrol"
	"fmt"
	"github.com/HairyMezican/Middleware/encapsulator"
	"github.com/HairyMezican/Middleware/errorhandler"
	"github.com/HairyMezican/Middleware/interceptor"
	"github.com/HairyMezican/Middleware/logger"
	"github.com/HairyMezican/Middleware/methoder"
	"github.com/HairyMezican/Middleware/oauther"
	"github.com/HairyMezican/Middleware/oauther/facebooker"
	"github.com/HairyMezican/Middleware/oauther/googleplusser"
	"github.com/HairyMezican/Middleware/parser"
	"github.com/HairyMezican/Middleware/sessioner"
	"github.com/HairyMezican/Middleware/statuser"
	"github.com/HairyMezican/Middleware/websocketer"
	"github.com/HairyMezican/TheRack/httper"
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

type TestDisplay string

func (this TestDisplay) Run(vars map[string]interface{}, next func()) {
	fmt.Println(this)
	next()
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

	ws := websocketer.New()
	ws.OnOpen(websocketcontrol.OpenUp)
	ws.OnClose(websocketcontrol.CloseDown)

	//set up the rack
	rackup := rack.New()
	rackup.Add(logger.Set(os.Stdout, "Log Test - ", log.LstdFlags))
	rackup.Add(requestlogger.M)
	rackup.Add(defaulter.V{"Layout": "base"})
	rackup.Add(encapsulator.AddLayout)
	rackup.Add(statuser.SetErrorLayout)
	if mode != debug {
		rackup.Add(errorhandler.ErrorHandler) //in debug version, it's more useful to just let it crash, so we can get more error information
	}
	rackup.Add(sessioner.Middleware)
	rackup.Add(login.Middleware)
	rackup.Add(ws)
	rackup.Add(parser.Form)
	rackup.Add(methoder.Override)
	rackup.Add(cept)
	rackup.Add(clashgetter.QueueGetter)
	rackup.Add(controllers.Root)

	//alert the user as to where we are running
	var site = ":80"
	if mode == debug {
		site = ":3000"
	}

	fmt.Print("\n\nStarting at localhost" + site + "!\n\n\n")
	conn := httper.HttpConnection(site)
	err := conn.Go(rackup)

	//We're ready to go!
	//run each request through the rack!
	if err != nil {
		fmt.Print("Error: " + err.Error())
	}
}
