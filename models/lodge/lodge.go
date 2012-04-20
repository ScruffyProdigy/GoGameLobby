package lodge

import (
	"../"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
)

type LodgeCollection struct {
	collection *mgo.Collection
}

var L = new(LodgeCollection)

func init() {
	models.RegisterModel(L)
}

type Lodge struct {
	Name   string
	Masons []string
}

func (this Lodge) Url() string {
	return "/lodges/" + this.Name
}

func (*LodgeCollection) CollectionName() string {
	return "lodges"
}

func (*LodgeCollection) RouteName() string {
	return "lodges"
}

func (*LodgeCollection) VarName() string {
	return "Lodge"
}

func (this *LodgeCollection) SetCollection(collection *mgo.Collection) {
	this.collection = collection
}

func (this *LodgeCollection) GetIndices() []mgo.Index {
	return []mgo.Index{
		{
			Key:    []string{"name"},
			Unique: true,
		},
		{
			Key: []string{"masons"},
		},
	}
}

func (this *LodgeCollection) LodgeFromName(s string) *Lodge {
	var result Lodge

	query := bson.M{"name": s}

	err := this.collection.Find(query).One(&result)
	if err != nil {
		return nil
	}

	return &result
}

func (this *LodgeCollection) LodgesFromMason(clashtag string) []Lodge {

	query := bson.M{"masons": clashtag}
	count, err := this.collection.Find(query).Count()
	if err != nil {
		return nil
	}

	result := make([]Lodge, count)

	err = this.collection.Find(query).All(&result)
	if err != nil {
		return nil
	}

	return result
}

func (this *LodgeCollection) AllLodges(out *[]Lodge) error {
	return this.collection.Find(bson.M{}).All(out)
}

func (this *LodgeCollection) Indexer(s string) interface{} {
	return this.LodgeFromName(s)
}

func (this *LodgeCollection) AddLodge(lodge *Lodge) error {
	return this.collection.Insert(lodge)
}
