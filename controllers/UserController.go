package controller

import "../routing"
import "fmt"
import "strconv"

type User struct {
	name   string
	points int
}

func init() {
	users := []User{{"Cole", 3}, {"Ryan", 2}}

	var indexer = func(s string) routes.Variable {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}
		return users[i]
	}

	rest := map[string]routes.HandlerFunc{
		"index": func(res routes.Response, req routes.Request, vars routes.VariableList) {
			page := "<html><head><title>Users</title></head><body><ul>"
			for _, user := range users {
				page += "<li>" + user.name + "</li>"
			}
			page += "</body></html>"

			fmt.Fprint(res, page)
		},
		"show": func(res routes.Response, req routes.Request, vars routes.VariableList) {
			page := "<html><head><title>User</title></head><body><h1>" + vars["user"].(User).name + "</h1><p>" + strconv.Itoa(vars["user"].(User).points) + "</p></body></html>"

			fmt.Fprint(res, page)
		},
	}

	routes.Root.AddRoute(routes.Resource("users", rest, "user", indexer))
}
