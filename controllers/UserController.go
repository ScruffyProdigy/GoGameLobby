package controller

import (
	"../log"
	"../routes"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"net/http"
	"strconv"
)

type User struct {
	ClashTag string
	Points   int
}

func (this User) Url() string {
	return "/users/" + this.ClashTag
}

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal("Please Launch Mongo before running this\n")
		panic(err)
	}

	collection := session.DB("test").C("users")

	var indexer = func(s string) interface{} {
		var result User

		query := bson.M{"clashtag": s}

		err = collection.Find(query).One(&result)
		if err != nil {
			return nil
		}

		return result
	}

	rest := map[string]routes.HandlerFunc{
		"index": func(res routes.Responder, req *http.Request, vars map[string]interface{}) {

			var users []User
			err = collection.Find(bson.M{}).All(&users)
			if err != nil {
				panic(err)
			}

			vars["Users"] = users
			vars["Title"] = "Users"
			vars["Layout"] = "base"

			res.Render("users/index")
		},
		"show": func(res routes.Responder, req *http.Request, vars map[string]interface{}) {
			user := vars["User"].(User)
			vars["Title"] = user.ClashTag
			vars["Layout"] = "base"
			res.Render("users/show")
		},
		"new": func(res routes.Responder, req *http.Request, vars map[string]interface{}) {
			vars["Title"] = "New User"
			vars["Layout"] = "base"
			res.Render("users/new")
		},
		"create": func(res routes.Responder, req *http.Request, vars map[string]interface{}) {
			err := req.ParseForm()
			if err != nil {
				panic(err)
			}

			var user User
			user.ClashTag = req.FormValue("User[ClashTag]")
			user.Points, err = strconv.Atoi(req.FormValue("User[Points]"))
			if err != nil {
				panic(err)
			}

			err = collection.Insert(&user)
			if err != nil {
				panic(err)
			}

			res.RedirectTo(user)
		},
	}

	userResource := routes.Resource("users", rest, "User", indexer)
	userResource.Collection.AddRoute(routes.Get("foo", func(res routes.Responder, req *http.Request, vars map[string]interface{}) {
		user := User{"Foo", 1}
		err := collection.Insert(&user)
		if err != nil {
			panic(err)
		}
		res.RedirectTo(user)
	}))

	routes.Root.AddRoute(userResource)
}
