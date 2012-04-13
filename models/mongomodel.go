package model

import (
	"../log"
	"launchpad.net/mgo"
)

var session *mgo.Session

type Collection interface {
	Indexer(string) interface{}
	CollectionName() string
	RouteName() string
	VarName() string
	SetCollection(*mgo.Collection)
}

var ModelList []Collection = make([]Collection, 0, 1)

func RegisterModel(model Collection) {
	ModelList = append(ModelList, model)
	model.SetCollection(session.DB("test").C(model.CollectionName()))
}

func init() {
	var err error
	session, err = mgo.Dial("localhost")
	if err != nil {
		log.Fatal("Please Launch Mongo before running this\n")
		panic(err)
	}
}
