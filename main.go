package main

import _ "./controllers"
import "./routing"

func main() {
	routes.Implement(routes.Root)
}
