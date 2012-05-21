package controllers

import (
	"../controller"
	"../login"
	"../models"
	"../models/game"
	"../models/lodge"
	"github.com/HairyMezican/Middleware/logger"
)

type ProjectController struct {
	controller.Heart
}

type modeDisplay struct {
	Identifier string
	Descriptor string
	Groups     []groupDisplay
	Groupable  bool
}

type groupDisplay struct {
	Identifier string
	Descriptor string
}

func (this modeDisplay) Playable() bool {
	return len(this.Groups) > 0
}

func (this modeDisplay) Multigrouped() bool {
	return len(this.Groups) > 1
}

type joinDisplay struct {
	Identifier string
	Descriptor string
}

func (ProjectController) RouteName() string {
	return "projects"
}

func (ProjectController) VarName() string {
	return "Game"
}

func (this ProjectController) Show() {
	log := (logger.V)(this.Vars).Get()
	g, isGame := this.Get("Game").(*game.Game)
	if !isGame {
		panic("Can't find Game")
	}

	this.Set("Title", g.Name)

	currentUser, loggedIn := (login.V)(this.Vars).CurrentUser()
	if g.Live && loggedIn {
		gameModes, err := g.GetGameModes(currentUser)
		if err != nil {
			this.Set("Error", "Could not contact host site:"+err.Error())
			this.Finish()
			return
		}

		// set up mode display
		var modes []modeDisplay
		for modeID, modeInfo := range gameModes {
			newMode := modeDisplay{
				Identifier: modeID,
				Descriptor: modeInfo.Name,
				Groupable:  true, //FIXME - should check to see if there is either more than one group, or if the one group has more than one player
				Groups:     []groupDisplay{},
			}

			for groupID, groupInfo := range *modeInfo.Groups {
				newGroup := groupDisplay{
					Identifier: groupID,
					Descriptor: groupInfo.Name,
				}

				newMode.Groups = append(newMode.Groups, newGroup)
			}

			modes = append(modes, newMode)
		}

		this.Set("GameModes", modes)
		log.Println(modes)
	}
}

func (this ProjectController) Update() {
	g, isGame := this.Get("Game").(*game.Game)
	if !isGame {
		panic("Can't find Game")
	}

	//TODO: Permissions

	lg := (logger.V)(this.Vars).Get()
	command := this.GetFormValue("Command")
	lg.Print("Command: " + command)

	switch this.GetFormValue("Command") {
	default:
		panic("Unrecognized Command")
	case "GoLive":
		g.Live = true
	case "NoLive":
		g.Live = false
		g.Modes = nil
	case "SetCommUrl":
		g.CommUrl = this.GetFormValue("CommUrl")
		g.Live = true
	}

	err := model.Save(g)
	if err != nil {
		errs, ok := err.(model.ValidationErrors)
		if !ok {
			panic(err)
		}

		for _, err := range errs {
			this.AddFlash(err.Field + " : " + err.Err)
		}
	}
	this.RedirectTo(g.Url())
}

func (this ProjectController) Create() {
	l, isLodge := this.Get("Lodge").(*lodge.Lodge)
	if !isLodge {
		panic("Cannot find lodge")
	}

	defer func() {
		rec := recover()
		if rec != nil {
			this.RedirectTo(l.Url())
			this.AddFlash("You fucked something up, please try again")
		}
	}()

	g := game.NewGame()
	g.Name = urlify(this.GetFormValue("Game[Name]"))
	g.Lodge = l.Name

	errs := model.Save(g)
	if errs != nil {
		panic(errs)
	}

	this.RespondWith(g)
}

func (this ProjectController) PostMemberJoin() {
	g, isGame := this.Get("Game").(*game.Game)
	if !isGame {
		panic("Can't find Game")
	}

	currentUser, loggedIn := (login.V)(this.Vars).CurrentUser()
	if !loggedIn {
		this.RedirectTo("/")
		return
	}
	mode := this.GetFormValue("mode")
	group := this.GetFormValue("group")
	join := this.GetFormValue("join")

	if join != "" {
		//we've got all the info we need to join
		this.join(mode, group, join)
		return
	}

	joinMethods, err := g.GetJoinModes(currentUser, mode, group)

	if err != nil || len(joinMethods) == 0 {
		this.AddFlash("There was a problem joining the game")
		this.RedirectTo(g.Url())
		return
	}

	if len(joinMethods) == 1 {
		//only one way to join, so might was well automatically do it
		for join, _ := range joinMethods {
			this.join(mode, group, join)
			return
		}
	}

	//display the user his joining choices
	var joins []joinDisplay
	for joinID, joinInfo := range joinMethods {
		newJoin := joinDisplay{
			Identifier: joinID,
			Descriptor: joinInfo.Name,
		}
		joins = append(joins, newJoin)
	}
	this.Set("Joins", joins)
	this.Render("join")
}

func (this ProjectController) join(mode, group, join string) {
	g, _ := this.Get("Game").(*game.Game)
	u, _ := (login.V)(this.Vars).CurrentUser()

	m, err := g.GetMode(mode)
	if err != nil {
		this.AddFlash("Unable to join the game right now")
	} else {
		m.AddToQueue(u.ClashTag, group, join)
	}

	this.RedirectTo(g.Url())
}

func init() {
	controller.RegisterController(&ProjectController{}, "projects", "Game", func(query string, vars map[string]interface{}) (interface{}, bool) {
		l, isLodge := vars["Lodge"].(*lodge.Lodge)
		if !isLodge {
			panic("Cannot find lodge")
		}

		result := game.G.GameFromLodgeAndName(l, query)
		return result, result != nil
	}).AddAsSubresource(Lodge)
}
