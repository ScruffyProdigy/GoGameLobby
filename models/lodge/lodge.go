package lodge

import (
	"../"
	"../user"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var L = new(LodgeCollection)

func init() {
	model.RegisterModel(L)
}

/*

	Lodge

*/

type Lodge struct {
	ID     bson.ObjectId `_id`
	Name   string
	Masons []string //could use ObjectID, but putting in the name allows us to not need to load the users or games into memory quite as often
	Games  []string //downside is, if any of the games or users gets renamed, we'll have to readjust all of the lodges that point to them
}

func NewLodge() *Lodge {
	return new(Lodge)
}

//		Utility Functions 

func (this Lodge) Url() string {
	return "/lodges/" + this.Name + "/"
}

func (this *Lodge) AddMason(u *user.User) {
	this.Masons = append(this.Masons, u.ClashTag)
	if !this.IsNew() {
		model.Save(this)
	}
	u.Lodges = append(u.Lodges, this.Name)
	if !u.IsNew() {
		model.Save(u)
	}
}

//		Interface Methods

func (this *Lodge) Validate() (errors *model.ValidationErrors) {
	//Name should be unique
	errors = model.NoErrors()

	other := L.LodgeFromName(this.Name)
	if other != nil && other.ID != this.ID {
		errors.Add("Name", "should be unique")
	}

	if this.Name == "new" {
		errors.Add("Name", "the word 'new' is reserved")
	}

	return
}

func (this *Lodge) IsNew() bool {
	return !this.ID.Valid()
}

func (this *Lodge) Collection() model.Collection {
	return L
}

func (this *Lodge) GetID() bson.ObjectId {
	return this.ID
}

func (this *Lodge) SetID(id bson.ObjectId) {
	this.ID = id
}

/*

	Lodge Collection

*/

type LodgeCollection struct {
	collection *mgo.Collection
}

//	Setup Functions

func (*LodgeCollection) CollectionName() string {
	return "lodges"
}

func (this *LodgeCollection) SetCollection(collection *mgo.Collection) {
	this.collection = collection
}

func (this *LodgeCollection) GetCollection() *mgo.Collection {
	return this.collection
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

//	Queries

func (this *LodgeCollection) LodgeFromName(s string) *Lodge {
	var result Lodge

	query := bson.M{"name": s}

	err := this.collection.Find(query).One(&result)
	if err != nil {
		return nil
	}

	return &result
}

func (this *LodgeCollection) LodgesFromMason(u *user.User) []Lodge {
	var result []Lodge

	query := bson.M{"masons": u.ClashTag}

	err := this.collection.Find(query).All(&result)
	if err != nil {
		return nil
	}

	return result
}

func (this *LodgeCollection) AllLodges(out *[]Lodge) error {
	return this.collection.Find(bson.M{}).All(out)
}
