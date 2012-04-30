package model

import (
	"errors"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"log"
)

var db *mgo.Database

type ValidationError struct {
	field string
	err   string
}

type ValidationErrors []ValidationError

func NoErrors() *ValidationErrors {
	result := make(ValidationErrors, 0)
	return &result
}

func (this ValidationErrors) Error() (result string) {
	result = "{"
	for _, err := range this {
		result += err.field + ":" + err.err
	}
	result += "}"
	return
}

func (this *ValidationErrors) Add(field, err string) {
	*this = append(*this, ValidationError{field: field, err: err})
}

type Model interface {
	IsNew() bool
	Collection() Collection
	Validate() *ValidationErrors
	GetID() bson.ObjectId
	SetID(bson.ObjectId)
}

func Save(m Model) error {
	errs := m.Validate()
	if errs != nil && len(*errs) > 0 {
		return errs
	}

	coll := m.Collection()
	if coll == nil {
		return errors.New("Cannot Get Collection")
	}
	c := coll.GetCollection()
	if c == nil {
		return errors.New("Cannot Get Collection")
	}
	var err error
	if m.IsNew() {
		m.SetID(bson.NewObjectId())
		err = c.Insert(m)
	} else {
		err = c.Update(bson.M{"_id": m.GetID()}, m)
	}
	if err != nil {
		return err
	}
	return nil
}

type Collection interface {
	CollectionName() string
	SetCollection(*mgo.Collection)
	GetCollection() *mgo.Collection
	GetIndices() []mgo.Index
}

var modelList map[string]Collection = make(map[string]Collection)

func RegisterModel(model Collection) {
	name := model.CollectionName()
	modelList[name] = model
	model.SetCollection(db.C(name))
}

func SetUp() {
	for name, model := range modelList {
		for _, index := range model.GetIndices() {
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
