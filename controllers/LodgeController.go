package controllers

import (
	"../controller"
	"../login"
	"../models"
	"../models/lodge"
)

type LodgeController struct {
	controller.Heart
}

func (LodgeController) RouteName() string {
	return "lodges"
}

func (LodgeController) VarName() string {
	return "Lodge"
}

func (this LodgeController) Indexer(s string) (interface{}, bool) {
	result := lodge.L.LodgeFromName(s)
	return result, result != nil
}

func (this LodgeController) Index() controller.Response {
	var lodges []lodge.Lodge
	err := lodge.L.AllLodges(&lodges)
	if err != nil {
		panic(err)
	}

	this.Set("Lodges", lodges)
	this.Set("Title", "Mason Lodges")

	return this.DefaultResponse()
}

func (this LodgeController) Show() controller.Response {
	l, isLodge := this.Get("Lodge").(*lodge.Lodge)
	if !isLodge {
		panic("Can't find lodge")
	}

	this.Set("Title", l.Name)

	return this.DefaultResponse()
}

func (this LodgeController) New() controller.Response {

	this.Set("Title", "Create a Mason Lodge")

	return this.DefaultResponse()
}

func (this LodgeController) Create() (response controller.Response) {
	defer func() {
		rec := recover()
		if rec != nil {
			errs, isSlice := rec.([]error)
			if isSlice {
				for _, err := range errs {
					this.AddFlash(err.Error())
				}
			}
			err, isError := rec.(error)
			if isError {
				this.AddFlash(err.Error())
			}
			str, isStr := rec.(string)
			if isStr {
				this.AddFlash(str)
			}
			this.AddFlash("You fucked something up, please try again")
			response = this.Redirection("/lodges/new")
		}
	}()

	l := lodge.NewLodge()

	l.Name = urlify(this.GetFormValue("Lodge[Name]"))
	user, loggedIn := login.CurrentUser(this.Vars)
	if !loggedIn {
		panic("you must be logged in to create a mason lodge")
	}
	l.AddMason(user)

	errs := model.Save(l)
	if errs != nil {
		panic(errs)
	}

	return this.RespondWith(l)
}

var Lodge *controller.ControllerShell

func init() {
	Lodge = controller.RegisterController(&LodgeController{})
	Lodge.AddToRoot()
}
