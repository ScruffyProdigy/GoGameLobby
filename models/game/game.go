package game

import (
	"../"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
)

var G = new(GameCollection)

func init() {
	model.RegisterModel(G)
}

/*

	Game

*/

type Game struct {
	ID bson.ObjectId `_id`
}

//		Interface Methods


func (this *Game) Validate() []error {
	return nil
}	

func (this *Game) IsNew() bool {
	return !this.ID.Valid()
}

func (*Game) Collection() model.Collection {
	return G
}

func (this *Game) GetID() bson.ObjectId {
	return this.ID
}

func (this *Game) SetID(id bson.ObjectId) {
	this.ID = id
}
//		Utility Functions


/*

	GameCollection

*/

type GameCollection struct {
	collection *mgo.Collection
}

//		Setup Functions


func (*GameCollection) CollectionName() string {
	return "games"
}

func (this *GameCollection) SetCollection(collection *mgo.Collection) {
	this.collection = collection
}

func (this *GameCollection) GetCollection() *mgo.Collection{
	return this.collection
}

func (this *GameCollection) GetIndices() []mgo.Index {
	return []mgo.Index{}
}


//		Queries


