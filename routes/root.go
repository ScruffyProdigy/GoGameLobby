package routes

//Root is the default Root to the directory structure
var Root RouteBranch

//if you want a directory structure separate from the Root, just call NewRoot to get one
func NewRoot() RouteBranch {
	return newRouteList()
}

func init() {
	Root = NewRoot()
}
