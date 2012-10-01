package user

import (
	"../"
	"../../global"
	"github.com/HairyMezican/SimpleRedis/redis"
	"github.com/HairyMezican/goauth2/oauth"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var U *UserCollection = new(UserCollection)

func init() {
	model.RegisterModel(U)
}

/*

	User

*/

type User struct {
	ID             bson.ObjectId `_id`
	ClashTag       string
	Points         int
	Authorizations []AuthorizationData
	Lodges         []string //we could use ObjectId's, but then we'd have to load the lodge into memory more often to get it's name
	//of course, by using the lodge's name, it means if the lodge ever gets renamed, we'll have to update all of the users
}

type AuthorizationData struct {
	Authorization string
	Id            string
	Token         oauth.Token
}

func NewUser() *User {
	return new(User)
}

//	Utility Functions

func (this User) Url() string {
	return "/users/" + this.ClashTag + "/"
}

//	Interface Methods

func (this *User) Validate() (errors *model.ValidationErrors) {
	errors = model.NoErrors()

	other := U.UserFromClashTag(this.ClashTag)
	if other != nil && other.ID != this.ID {
		errors.Add("ClashTag", "should be unique")
	}

	if this.ClashTag == "new" {
		errors.Add("ClashTag", "the word 'new' is reserved")
	}

	return
}

func (this *User) IsNew() bool {
	return !this.ID.Valid()
}

func (this *User) Collection() model.Collection {
	return U
}

func (this *User) GetID() bson.ObjectId {
	return this.ID
}

func (this *User) SetID(id bson.ObjectId) {
	this.ID = id
}

func (this *User) Queues() redis.Set {
	return global.Redis.Set("users " + this.ClashTag + " queues")
}

func (this *User) Clashes() redis.Set {
	return global.Redis.Set("users " + this.ClashTag + " clashes")
}

/*

	UserCollection

*/

type UserCollection struct {
	collection *mgo.Collection
}

//	Setup Functions

func (*UserCollection) CollectionName() string {
	return "users"
}

func (this *UserCollection) SetCollection(collection *mgo.Collection) {
	this.collection = collection
}

func (this *UserCollection) GetCollection() *mgo.Collection {
	return this.collection
}

func (this *UserCollection) GetIndices() []mgo.Index {
	return []mgo.Index{
		{
			Key:    []string{"clashtag"},
			Unique: true,
		},
		{
			Key:    []string{"authorizations.authorization", "authorizations.id"},
			Unique: true,
		},
	}
}

//	Queries

func (this *UserCollection) UserFromClashTag(s string) *User {
	var result User

	query := bson.M{"clashtag": s}

	err := this.collection.Find(query).One(&result)
	if err != nil {
		return nil
	}

	return &result
}

func (this *UserCollection) UserFromAuthorization(authorization, id string) *User {
	var result User

	query := bson.M{"authorizations.authorization": authorization, "authorizations.id": id}

	err := this.collection.Find(query).One(&result)
	if err != nil {
		return nil
	}

	return &result
}

func (this *UserCollection) AllUsers(out *[]User) error {
	return this.collection.Find(bson.M{}).All(out)
}
