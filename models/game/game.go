package game

import (
	"../"
	"../../gamedata"
	"../../messenger"
	"../../redis"
	"../lodge"
	"../user"
	"errors"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
)

const (
	sEp = " "
)

var G = new(GameCollection)
var QueueMutex *redis.ReadWriteMutex

func init() {
	model.RegisterModel(G)
	var err error
	QueueMutex, err = redis.RWMutex("QueueMutex", 16)
	if err != nil {
		panic("Couldn't set up Queue Mutex")
	}
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
	Modes   map[string]Mode
}

type Mode struct {
	GroupCount map[string]int //a list of the groups needed for the mode, and the number of people needed to fill the group
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
	return messenger.JSONmessage(message, this.CommUrl, result)
}

func (this *Game) GetPlayerCounts(modename string) (gamedata.ModePlayerCounts, error) {
	var Data struct {
		Mode gamedata.PrePlayerCountInfo `json:"playercount"`
	}
	Data.Mode.Mode = modename

	var Result gamedata.ModePlayerCounts
	err := this.Message(Data, &Result)

	return Result, err
}

func (this *Game) GetGameModes(u *user.User) (map[string]gamedata.ModeInfo, error) {
	var Data struct {
		Mode gamedata.PreModeInfo `json:"modeinfo"`
	}
	Data.Mode.User = u.ClashTag

	var Result map[string]gamedata.ModeInfo
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
	err := this.Message(Data, &Result)

	return Result, err
}

func (this *Game) StartClash(mode string, PreStart map[string]map[string]string) (gamedata.StartInfo, error) {
	var Data struct {
		Start gamedata.PreStartInfo `json:"startinfo"`
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
	err := this.Message(Data, &Result)

	return Result, err
}

func (this *Game) getMode(mode string) Mode {
	m, ok := this.Modes[mode]
	if !ok {
		m = this.SetUpMode(mode)
	}
	return m
}

func (this *Game) startMode(mode string) (restart bool) {
	game := this.Name
	m := this.getMode(mode)

	var startInfo gamedata.StartInfo
	var err error
	start := QueueMutex.Write.Try(func() {
		preStartInfo := make(map[string]map[string]string)
		for group, full := range m.GroupCount {
			preStartInfo[group] = make(map[string]string)
			for i := 0; i < full; i++ {
				user := pullFromQueue(game, mode, group)
				join := takeJoinOptions(user, game, mode)
				preStartInfo[group][user] = join
				removeFromMode(game, mode, user)
			}
		}

		startInfo, err = this.StartClash(mode, preStartInfo)
		if err != nil {
			//there was an error trying to start the clash, put everybody back in the queue
			for group, players := range preStartInfo {
				for user, join := range players {
					addToMode(user, game, mode)
					jumpTheQueue(user, game, mode, group)
					setJoinOptions(user, game, mode, join)
				}
			}
		}
	})

	if err != nil {
		return
	}

	if start {
		// we've loaded the players in the queue, time to start it
		sendStartMessages(startInfo)
		return false
	}

	// somebody else was trying to start a game, loop through again to check to see if the game can start
	return true
}

func (this *Game) checkForStart(mode string) {
	game := this.Name
	m := this.getMode(mode)

	start := true

	trying := true

	for trying {
		//if the queue is long enough to start, start it
		QueueMutex.Read.Force(func() {
			for group, full := range m.GroupCount {
				count := int(queueLength(game, mode, group))
				if count < full {
					start = false
					return
				}
			}
		})

		if !start {
			return
		}

		trying = this.startMode(mode)
	}

}

func (this *Game) AddToQueue(user, mode, group, join string) (err error) {

	defer func() {
		rec := recover()
		if rec != nil {
			if e, ok := rec.(error); ok {
				err = e
				return
			}
			if str, ok := rec.(string); ok {
				err = errors.New(str)
				return
			}
			err = errors.New("UnknownError")
			return
		}
	}()

	game := this.Name

	//add the user to the set of players queuing for a clash (and find out if they're already on the list)
	isNew := addToMode(user, this.Name, mode)
	if !isNew {
		m := this.getMode(mode)
		//the player is already part of the game, remove them from their old queue first
		for group, _ := range m.GroupCount {
			if removeFromQueue(user, game, mode, group) {
				break
			}
		}
	}

	//add the user to the queue
	addToQueue(user, game, mode, group)
	setJoinOptions(user, game, mode, join)

	go this.checkForStart(mode)

	return nil
}

func (this *Game) RemoveFromQueue(user, mode string) (err error) {
	defer func() {
		rec := recover()
		if rec != nil {
			if e, ok := rec.(error); ok {
				err = e
				return
			}
			if str, ok := rec.(string); ok {
				err = errors.New(str)
				return
			}
			err = errors.New("UnknownError")
			return
		}
	}()

	m := this.getMode(mode)
	game := this.Name

	if !removeFromMode(user, game, mode) {
		//not in this mode
		return
	}
	for group, _ := range m.GroupCount {
		if removeFromQueue(user, game, mode, group) {
			break
		}
	}
	takeJoinOptions(user, game, mode)

	return nil
}

func (this *Game) SetUpMode(mode string) Mode {
	playercounts, err := this.GetPlayerCounts(mode)
	if err != nil {
		panic(err)
	}

	if this.Modes == nil {
		this.Modes = make(map[string]Mode)
	}
	this.Modes[mode] = Mode{
		GroupCount: playercounts.Players,
	}

	model.Save(this)

	return this.Modes[mode]
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
