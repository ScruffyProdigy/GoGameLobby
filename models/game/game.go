package game

import (
	"../"
	"../../gamedata"
	"../../messenger"
	"../lodge"
	"../user"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	sEp = " "
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
	Name      string
	Published bool
	Live      bool
	Lodge     string //could use ObjectID, but putting in the name allows us to not need to load the lodge into memory quite as often
	//downside is, if the lodge gets renamed, we'll have to readjust all of the games that point to it
	CommUrl string
	Modes   map[string]*Mode
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
		return this.GameUrl()
	}
	return this.ProjectUrl()
}

func (this *Game) ProjectUrl() string {
	return "/lodges/" + this.Lodge + "/projects/" + this.Name + "/"
}

func (this *Game) GameUrl() string {
	return "/games/" + this.Name + "/"
}

func (this *Game) Message(message, result interface{}) error {
	print("Querying Game - ")
	err := messenger.JSONmessage(message, this.CommUrl, result)
	print("Response Received\n")
	return err
}

func (this *Game) GetPlayerCounts(modename string) (gamedata.ModePlayerCounts, error) {
	var Data struct {
		Mode gamedata.PrePlayerCountInfo `json:"playercount"`
	}
	Data.Mode.Mode = modename

	var Result gamedata.ModePlayerCounts
	print("Need To Get Player Counts\n")
	err := this.Message(Data, &Result)

	return Result, err
}

func (this *Game) GetGameModes(u *user.User) (map[string]gamedata.ModeInfo, error) {
	var Data struct {
		Mode gamedata.PreModeInfo `json:"modeinfo"`
	}
	Data.Mode.User = u.ClashTag

	var Result map[string]gamedata.ModeInfo
	print("Need to Get Game Modes\n")
	err := this.Message(Data, &Result)

	return Result, err
}

func (this *Game) GetJoinModes(u *user.User, mode string, group string) (map[string]gamedata.JoinInfo, error) {
	var Data struct {
		Join gamedata.PreJoinInfo `json:"joininfo"`
	}
	Data.Join.User = u.ClashTag
	Data.Join.Mode = mode
	Data.Join.Group = group

	var Result map[string]gamedata.JoinInfo
	print("Need to Get Join Modes\n")
	err := this.Message(Data, &Result)

	return Result, err
}

func (this *Game) StartClash(mode string, PreStart map[string]map[string]string) (gamedata.StartInfo, error) {
	var Data struct {
		Start gamedata.PreStartInfo `json:"start"`
	}

	Data.Start.Mode = mode
	Data.Start.Groups = make(map[string]gamedata.GroupJoinInfo)

	for group, playeroptions := range PreStart {
		groupoptions := gamedata.GroupJoinInfo{}
		groupoptions.Players = make(map[string]gamedata.PlayerJoinInfo)

		Data.Start.Groups[group] = groupoptions

		for player, options := range playeroptions {
			join := gamedata.PlayerJoinInfo{}
			join.Join = options

			Data.Start.Groups[group].Players[player] = join
		}
	}

	var Result gamedata.StartInfo
	print("Time to Start the Clash!\n")
	err := this.Message(Data, &Result)

	return Result, err
}

func (this *Game) GetMode(mode string) (*Mode, error) {
	m, ok := this.Modes[mode]
	if !ok {
		var err error
		m, err = this.SetUpMode(mode)
		if err != nil {
			return nil, err
		}
	}
	m.game = this.Name
	m.mode = mode
	return m, nil
}

func (this *Game) SetUpMode(mode string) (*Mode, error) {
	playercounts, err := this.GetPlayerCounts(mode)
	if err != nil {
		return nil, err
	}

	if this.Modes == nil {
		this.Modes = make(map[string]*Mode)
	}

	m := new(Mode)
	m.GroupCount = playercounts.Players

	this.Modes[mode] = m

	model.Save(this)

	return m, nil
}

func (this *Game) ClearModes() {
	this.Modes = nil
	model.Save(this)
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
			Key: []string{"lodge", "published"},
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

func (this *GameCollection) GameFromLodgeAndName(l *lodge.Lodge, name string) *Game {
	var result Game

	query := bson.M{"lodge": l.Name, "name": name}

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

func (this *GameCollection) UnpublishedGamesFromLodge(l *lodge.Lodge) []Game {
	var result []Game

	query := bson.M{"lodge": l.Name, "published": false}

	err := this.collection.Find(query).All(&result)
	if err != nil {
		return nil
	}

	return result
}
