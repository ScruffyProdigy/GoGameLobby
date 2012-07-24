package controllers

import (
	"github.com/HairyMezican/TheRack/httper"
	"../controller"
	"../models/clash"
	"encoding/json"
	"../login"
)

type ClashController struct {
	controller.Heart
}

func (this ClashController) Get() {
	c := this.GetVar("Clash").(clash.Clash)
	
	//find out if the current user is a player in the chosen clash
	user, loggedIn := (login.V)(this.Vars).CurrentUser()
	
	//redirect the user to the clash (if they are a participant, set the url vals accordingly)
	if loggedIn {
		this.RedirectTo(c.PlayerUrl(user.ClashTag))
	} else {
		this.RedirectTo(c.Url())
	}
}

func (this ClashController) Update() {
	c := this.GetVar("Clash").(clash.Clash)
	r := (httper.V)(this.Vars).GetRequest()
	var m struct {
		Remove *string `json:"remove"`
		Results *[][]string `json:"results"`
	}
	json.NewDecoder(r.Body).Decode(&m)
	//either removing a player
	if m.Remove != nil {
		c.RemovePlayer(*m.Remove)
	}
	//or posting the results
	if m.Results != nil {
		c.Results(*m.Results)
	}
}

func init() {
	controller.RegisterController(&ClashController{}, "clashes", "Clash", func(s string, vars map[string]interface{}) (interface{}, bool) {
		c := clash.FromHash(clash.Hash(s))
		if !c.Exists() {
			return nil,false
		}
		return c,true
	}).AddTo(Root)
}