package controller

import (
	"../log"
	"../routing"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"net/http"
)

type User struct {
	GamerTag string
	Points   int
}

func (this User) Url() []string {
	return []string{"users", this.GamerTag}
}

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Error("\nError - Please Launch Mongo before running this\n")
		panic(err)
	}

	collection := session.DB("test").C("users")

	var indexer = func(s string) interface{} {
		var result User

		query := bson.M{"gamertag": s}

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

			vars["users"] = users
			vars["title"] = "Users"
			vars["layout"] = "base"

			res.Render("users/index")
		},
		"show": func(res routes.Responder, req *http.Request, vars map[string]interface{}) {
			user := vars["user"].(User)
			vars["title"] = user.GamerTag
			vars["layout"] = "base"
			res.Render("users/show")
		},
	}

	userResource := routes.Resource("users", rest, "user", indexer)
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
