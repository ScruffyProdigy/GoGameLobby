package models

import (
	"../log"
	"launchpad.net/mgo"
)

var db *mgo.Database

type Collection interface {
	GetIndices() []mgo.Index
	Indexer(string) interface{}
	CollectionName() string
	RouteName() string
	VarName() string
	SetCollection(*mgo.Collection)
}

var modelList map[string]Collection = make(map[string]Collection)

func RegisterModel(model Collection) {
	name := model.CollectionName()
	modelList[name] = model
	model.SetCollection(db.C(name))
}

func SetUp() {
	for name,model := range(modelList) {
		for _,index := range(model.GetIndices()) {
			db.C(name).EnsureIndex(index)
		}
	}
}

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal("Please Launch Mongo before running this\n")
		panic(err)
	}
	db = session.DB("test")
}
