package main

import _ "./controllers"
import "./routing"
import "./rack"
import "./templater"

func main() {
	templater.LoadTemplates("./views")
	rack.Up.Add(rack.Parser)
	rack.Up.Add(routes.RouteWare(routes.Root))
	rack.Up.Add(routes.RenderWare)
	rack.Up.Go(rack.HttpConnection(":3000"))
}
