package user

import (
	"../"
	"../../goauth2/oauth"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
)

type UserCollection struct {
	collection *mgo.Collection
}

var U *UserCollection = new(UserCollection)

func init() {
	model.RegisterModel(U)
}

type FacebookUserData struct {
	Id    int
	Token oauth.Token
}

type User struct {
	ClashTag string
	Points   int
	Facebook []FacebookUserData
}

func (this User) Url() string {
	return "/users/" + this.ClashTag
}

func (*UserCollection) CollectionName() string {
	return "users"
}

func (*UserCollection) RouteName() string {
	return "users"
}

func (*UserCollection) VarName() string {
	return "User"
}

func (this *UserCollection) SetCollection(collection *mgo.Collection) {
	this.collection = collection
}

func (this *UserCollection) UserFromClashTag(s string) *User {
	var result User

	query := bson.M{"clashtag": s}

	err := this.collection.Find(query).One(&result)
	if err != nil {
		return nil
	}

	return &result
}

/*
	var indexer = func(s string) interface{} {
		var result User

		query := bson.M{"clashtag": s}

		err = collection.Find(query).One(&result)
		if err != nil {
			return nil
		}

		return result
	}
*/
func (this *UserCollection) UserFromFacebookID(id int) *User {
	var result User

	query := bson.M{"facebook.id": id}

	err := this.collection.Find(query).One(&result)
	if err != nil {
		return nil
	}

	return &result
}

func (this *UserCollection) AllUsers(out *[]User) error {
	return this.collection.Find(bson.M{}).All(out)
}

func (this *UserCollection) Indexer(s string) interface{} {
	return this.UserFromClashTag(s)
}

func (this *UserCollection) AddUser(user *User) error {
	return this.collection.Insert(user)
}
