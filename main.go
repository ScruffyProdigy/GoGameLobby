package main

import (
	_ "./controllers"
	"./layouts"
	"./rack"
	"./routes"
	"./session"
	"./templater"
)

func main() {
	templater.LoadTemplates("./views")

	rack.Up.Add(layouts.AddLayout)
	rack.Up.Add(layouts.SetErrorLayout)
	rack.Up.Add(layouts.Defaulter("Layout", "base"))
	rack.Up.Add(session.Middleware)
	rack.Up.Add(routes.Parser)
	rack.Up.Add(routes.RouterWare(routes.Root))

	rack.Up.Go(rack.HttpConnection(":3000"))
}
