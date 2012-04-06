package controller

import (
	"../log"
	"../routing"
	"fmt"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"net/http"
	"strconv"
)

type User struct {
	GamerTag string
	Points   int
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
		"index": func(res http.ResponseWriter, req *http.Request, vars map[string]interface{}) {

			var users []User
			err = collection.Find(bson.M{}).All(&users)
			if err != nil {
				panic(err)
			}

			page := "<html><head><title>Users</title></head><body><ul>"
			for _, user := range users {
				page += "<li><a href='/users/" + user.GamerTag + "'>" + user.GamerTag + " - " + strconv.Itoa(user.Points) + "</a></li>"
			}
			page += "</body></html>"

			fmt.Fprint(res, page)
		},
		"show": func(res http.ResponseWriter, req *http.Request, vars map[string]interface{}) {
			page := "<html><head><title>User</title></head><body><h1>" + vars["user"].(User).GamerTag + "</h1><p>" + strconv.Itoa(vars["user"].(User).Points) + "</p></body></html>"

			fmt.Fprint(res, page)
		},
	}

	userResource := routes.Resource("users", rest, "user", indexer)
	userResource.Collection.AddRoute(routes.Get("foo", func(res http.ResponseWriter, req *http.Request, vars map[string]interface{}) {
		err := collection.Insert(&User{"Foo", 1})
		if err != nil {
			panic(err)
		}
		fmt.Fprint(res, "Okay, it's done!")
	}))

	routes.Root.AddRoute(userResource)
}
