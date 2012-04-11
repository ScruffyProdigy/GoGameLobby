package main

import (
	_ "./controllers"
	"./layouts"
	"./rack"
	"./routes"
	"./session"
	"./templater"
	"./errorhandler"
	"./notfound"
)

func main() {
	templater.LoadTemplates("./views")

	rack.Up.Add(layouts.AddLayout)
	rack.Up.Add(layouts.SetErrorLayout)
	rack.Up.Add(layouts.Defaulter("Layout", "base"))
	rack.Up.Add(errorhandler.ErrorHandler)
	rack.Up.Add(session.Middleware)
	rack.Up.Add(routes.Parser)
	rack.Up.Add(routes.RouterWare(routes.Root))
	rack.Up.Add(notfound.NotFound)

	rack.Up.Go(rack.HttpConnection(":3000"))
}
