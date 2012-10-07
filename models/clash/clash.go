package clash

import (
	"../../gamedata"
	"../../global"
	"../../pubsuber"
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"github.com/HairyMezican/SimpleRedis/redis"
	neturl "net/url"
)

type Hash string

const (
	hashLength = 128
)

func NewHash() Hash {
	var randomBytes [hashLength]byte
	rand.Read(randomBytes[:])

	en := base32.StdEncoding
	result := make([]byte, en.EncodedLen(len(randomBytes)))
	en.Encode(result, randomBytes[:])

	return Hash(bytes.ToLower(result))
}

type Clash struct {
	id Hash
}

/*
Private methods
*/

//Redis Functions

func (this *Clash) Prefix() redis.Prefix {
	return global.Redis.Prefix("Clash " + string(this.id) + " ")
}

func (this *Clash) base() redis.Hash {
	return this.Prefix().Hash("Base")
}

func (this *Clash) url() redis.HashString {
	return this.base().String("Url")
}

func (this *Clash) playerUrl(clashtag string) redis.HashString {
	return this.base().String("PlayerUrl " + clashtag)
}

func (this *Clash) game() redis.HashString {
	return this.base().String("Game")
}

func (this *Clash) mode() redis.HashString {
	return this.base().String("Mode")
}

func (this *Clash) players() redis.Set {
	return this.Prefix().Set("Players")
}

func userClashes(clashtag string) redis.Set {
	//return user.Prefix().Set("Clashes")
	return global.Redis.Prefix("User " + clashtag + " ").Set("Clashes")
}

//constructor&destructor
func (this *Clash) addPlayer(player string, options gamedata.PlayerStartInfo) {
	this.players().Add(player)
	userClashes(player).Add(string(this.id))

	urlvals := make(neturl.Values)
	for k, v := range options.UrlValues {
		urlvals.Add(k, v)
	}
	query := urlvals.Encode()

	//TODO: We had the url, game and name in the function before, come up with a more elegant way of transfering them
	url := <-this.url().Get()
	if query != "" {
		url += "?" + query
	}

	this.playerUrl(player).Set(url)
	pubsuber.User(player).SendMessage("start", map[string]string{"Game": this.Game(), "Mode": this.Mode()})
}

//Hack: I would just use mode directly, but it causes an import loop, so just using an interface of it
type Mode interface {
	Game() string
	Mode() string
}

func (this *Clash) Setup(game, mode string, psi gamedata.StartInfo) {
	this.url().Set(psi.Url)
	this.game().Set(game)
	this.mode().Set(mode)
	for player, options := range psi.Players {
		this.addPlayer(player, options)
	}
}

func (this *Clash) teardown() {
	this.base().Delete()
	this.players().Delete()
}

/*
Public Methods
*/

//	Getters
func New() *Clash {
	this := new(Clash)
	this.id = NewHash()
	return this
}

func FromUser(clashtag string) []Clash {
	clashhashes := <-userClashes(clashtag).Members()
	clashes := make([]Clash, len(clashhashes))
	for i, hash := range clashhashes {
		clashes[i].id = Hash(hash)
	}
	return clashes
}

func FromHash(h Hash) *Clash {
	c := new(Clash)
	c.id = h
	return c
}

//returns the clashtags of all involved users
func (this *Clash) Players() []string {
	return <-this.players().Members()
}

//returns the url necessary to go to this clash for the user
func (this *Clash) PlayerUrl(clashtag string) string {
	return <-this.playerUrl(clashtag).Get()
}

//removes a player from the clash
func (this *Clash) RemovePlayer(clashtag string) {
	userClashes(clashtag).Remove(string(this.id))
}

//finishes a clash, and lets it know what the results were
func (this *Clash) Results(results [][]string) {
	// TODO: update players scores here

	//delete the clash
	this.teardown()
}

//
func (this *Clash) Exists() bool {
	return <-this.base().Exists()
}

func (this *Clash) Hash() Hash {
	return this.id
}

func (this *Clash) Mode() string {
	return <-this.mode().Get()
}

func (this *Clash) Game() string {
	return <-this.game().Get()
}

func (this *Clash) Details(player string) (string, string, string) {
	playerUrlChan := this.playerUrl(player).Get()
	game := this.game().Get()
	mode := this.mode().Get()
	url, ok := <-playerUrlChan
	if ok {
		return <-game, <-mode, url
	}
	urlChan := this.url().Get()
	return <-game, <-mode, <-urlChan
}

func (this *Clash) Url() string {
	return <-this.url().Get()
}

func (this *Clash) ResponseUrl() string {
	return "/clashes/" + string(this.id) + "/"
}
