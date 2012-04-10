package main

import (
	_ "./controllers"
	"./middleware"
	"./rack"
	"./routes"
	"./session"
	"./templater"
)

func main() {
	templater.LoadTemplates("./views")

	rack.Up.Add(middleware.Parser)
	rack.Up.Add(middleware.AddLayout)
	rack.Up.Add(middleware.SetErrorLayout)
	rack.Up.Add(session.Middleware)
	rack.Up.Add(routes.RouterWare(routes.Root))

	rack.Up.Go(rack.HttpConnection(":3000"))
}
