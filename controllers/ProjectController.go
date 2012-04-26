package controllers

import (
	"../controller"
	"../models"
	"../models/game"
	"../models/lodge"
)

func init() {
	controller.RegisterController(&ProjectController{g: game.G}).AddAsSubresource(Lodge)
}

type ProjectController struct {
	g *game.GameCollection
	controller.Heart
}

func (this ProjectController) Indexer(query string) (interface{}, bool) {
	l, isLodge := this.Get("Lodge").(*lodge.Lodge)
	if !isLodge {
		panic("Cannot find lodge")
	}

	result := this.g.GameFromLodgeAndName(l.Name, query)
	return result, result != nil
}

func (ProjectController) RouteName() string {
	return "projects"
}

func (ProjectController) VarName() string {
	return "Game"
}

func (this ProjectController) Show() controller.Response {
	g, isGame := this.Get("Game").(*game.Game)
	if !isGame {
		panic("Can't find Game")
	}

	this.Set("Title", g.Name)
	return this.DefaultResponse()
}

func (this ProjectController) Create() (response controller.Response) {
	l, isLodge := this.Get("Lodge").(*lodge.Lodge)
	if !isLodge {
		panic("Cannot find lodge")
	}

	defer func() {
		rec := recover()
		if rec != nil {
			response = this.Redirection(l.Url())
			this.AddFlash("You fucked something up, please try again")

		}
	}()

	g := game.NewGame()
	g.Name = this.GetFormValue("Game[Name]")
	g.Lodge = l.Name

	errs := model.Save(g)
	if errs != nil {
		panic(errs)
	}

	return this.RespondWith(g)
}
