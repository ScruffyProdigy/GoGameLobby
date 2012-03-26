package main

import _ "./controllers"
import "./routing"
import "./rack"

func main() {
	rack.Up.Add(rack.Parser)
	rack.Up.Add(routes.EndWare(routes.Root))
	rack.Up.Go(rack.HttpConnection(":3000"))
}
