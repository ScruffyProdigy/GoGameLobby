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
	ID        bson.ObjectId `_id`
	Name      string        `name`
	Published bool          `published`
	Lodge     string        `lodge`
}

func NewGame() *Game {
	return new(Game)
}

//		Interface Methods

func (this *Game) Validate() (errors *model.ValidationErrors) {
	errors = model.NoErrors()

	other := G.GameFromName(this.Name)
	if other != nil && other.ID != this.ID {
		errors.Add("Name", "should be unique")
	}

	return
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
func (this *Game) Url() string {
	if this.Published {
		return "/games/" + this.Name + "/"
	}
	return "/lodges/" + this.Lodge + "/projects/" + this.Name + "/"
}

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

func (this *GameCollection) GetCollection() *mgo.Collection {
	return this.collection
}

func (this *GameCollection) GetIndices() []mgo.Index {
	return []mgo.Index{
		{
			Key:    []string{"name"},
			Unique: true,
		},
		{
			Key:    []string{"published", "name"},
			Unique: true,
		},
		{
			Key:    []string{"lodge", "name"},
			Unique: true,
		},
	}
}

//		Queries

func (this *GameCollection) GameFromName(name string) *Game {
	var result Game

	query := bson.M{"name": name}

	err := this.collection.Find(query).One(&result)
	if err != nil {
		return nil
	}

	return &result
}

func (this *GameCollection) GameFromLodgeAndName(lodge, name string) *Game {
	var result Game

	query := bson.M{"lodge": lodge, "name": name}

	err := this.collection.Find(query).One(&result)
	if err != nil {
		return nil
	}

	return &result
}

func (this *GameCollection) PublishedGameFromName(name string) *Game {
	var result Game

	query := bson.M{"published": true, "name": name}

	err := this.collection.Find(query).One(&result)
	if err != nil {
		return nil
	}

	return &result
}
