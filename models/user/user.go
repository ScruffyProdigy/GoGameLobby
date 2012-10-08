package user

import (
	"../"
	"../../global"
	"errors"
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
	//TODO: should probably change these to some other data structure than an unsorted list
	Lodges         []string //we could use ObjectId's, but then we'd have to load the lodge into memory more often to get it's name
	Friends        []string //once again, we could use ObjectId's but we'd have to load every user into memory to get their names
	FriendRequests []string //ditto
	//of course, by using the lodges' and the users' names, it means if they ever get renamed, we'll have to update all of the users
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

func (this *User) RequestFriend(friend string) error {
	if this.IsFriend(friend) {
		return errors.New("You are already Friends")
	}
	if err := this.AcceptRequest(friend); err == nil {
		return nil
	}
	other := U.UserFromClashTag(friend)
	other.FriendRequests = append(other.FriendRequests, this.ClashTag)
	model.Save(other)
	return nil
}

func (this *User) removeRequest(friend string) bool {
	for i, name := range this.FriendRequests {
		if name == friend {
			this.FriendRequests[i] = this.FriendRequests[len(this.FriendRequests)-1]
			this.FriendRequests = this.FriendRequests[:len(this.FriendRequests)-1]
			return true
		}
	}
	return false
}

func (this *User) HasRequest(friend string) bool {
	for _, name := range this.FriendRequests {
		if name == friend {
			return true
		}
	}
	return false
}

func (this *User) AcceptRequest(friend string) error {
	if this.removeRequest(friend) {
		this.AddFriend(friend)
		return nil
	}
	return errors.New("No Request Found")
}

func (this *User) DenyRequest(friend string) error {
	if this.removeRequest(friend) {
		model.Save(this)
		return nil
	}
	return errors.New("No Request Found")
}

func (this *User) IsFriend(friend string) bool {
	for _, f := range this.Friends {
		if f == friend {
			return true
		}
	}
	return false
}

func (this *User) addFriend(friend string) {
	this.Friends = append(this.Friends, friend)
	model.Save(this)
}

func (this *User) removeFriend(friend string) error {
	for i, f := range this.Friends {
		if f == friend {
			this.Friends[i] = this.Friends[len(this.Friends)-1]
			this.Friends = this.Friends[:len(this.Friends)-1]
			model.Save(this)
			return nil
		}
	}
	return errors.New("You aren't friends")
}

func (this *User) AddFriend(friend string) {
	this.addFriend(friend)
	U.UserFromClashTag(friend).addFriend(this.ClashTag)
}

func (this *User) Unfriend(friend string) {
	this.removeFriend(friend)
	U.UserFromClashTag(friend).removeFriend(this.ClashTag)
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
