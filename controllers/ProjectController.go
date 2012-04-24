package controller

import (
	"../log"
	"../models"
	"../models/game"
	"../models/lodge"
	"../rack"
	"../redirecter"
	"../routes"
	"../session"
	"fmt"
	"net/http"
)

func init() {
	routes.Resource(ProjectController{game.G}).AddTo(LodgeRoute.Member)
}

type ProjectController struct {
	G *game.GameCollection
}

func (this ProjectController) Indexer(query string, vars rack.Vars) (interface{}, bool) {
	l, isLodge := vars["Lodge"].(*lodge.Lodge)
	if !isLodge {
		panic("Cannot find lodge")
	}
	result := this.G.GameFromLodgeAndName(l.Name, query)
	return result, result != nil
}

func (ProjectController) RouteName() string {
	return "projects"
}

func (ProjectController) VarName() string {
	return "Game"
}

func (this ProjectController) Show(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
	g, isGame := vars["Game"].(*game.Game)
	if !isGame {
		panic("Can't find Game")
	}

	vars["Title"] = g.Name
	return next()
}

func (this ProjectController) Create(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	l, isLodge := vars["Lodge"].(*lodge.Lodge)
	if !isLodge {
		panic("Cannot find lodge")
	}

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Fprint(log.DebugLog(), "\n\nDebug wtf is rec?:", rec)
			errs, isSlice := rec.([]error)
			err, isError := rec.(error)
			str, isString := rec.(string)
			if isSlice {
				for _, err = range errs {
					session.AddFlash(err.Error())
				}
			} else if isError {
				session.AddFlash(err.Error())
			} else if isString {
				session.AddFlash(str)
			} else {
				session.AddFlash("unknown error")
			}
			status, header, message = redirecter.Go(r, vars, l.Url(),
				session.AddFlash("You fucked something up, please try again"))

		}
	}()

	var g game.Game
	g.Name = r.FormValue("Game[Name]")
	g.Published = false
	g.Lodge = l.Name

	errs := model.Save(&g)
	if errs != nil {
		panic(errs)
	}

	vars["Game"] = &g
	return next()
}
