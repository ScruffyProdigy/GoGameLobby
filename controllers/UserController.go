package controllers

import (
	"../controller"
	"../login"
	"../models"
	"../models/user"
	"../session"
	"strconv"
	"time"
)

type UserController struct {
	u *user.UserCollection
	controller.Heart
}

func (UserController) RouteName() string {
	return "users"
}

func (UserController) VarName() string {
	return "User"
}

func (this UserController) Indexer(s string) (interface{}, bool) {
	result := this.u.UserFromClashTag(s)
	return result, result != nil
}

func (this UserController) Index() controller.Response {
	var users []user.User
	err := this.u.AllUsers(&users)
	if err != nil {
		panic(err)
	}

	this.Set("Users", users)
	this.Set("Title", "Users")

	return this.DefaultResponse()
}

func (this UserController) Show() controller.Response {
	u := this.Get("User").(*user.User)

	this.Set("Title", u.ClashTag)

	return this.DefaultResponse()
}

func (this UserController) New() controller.Response {

	authorization, isString := this.Session().Clear("authorization").(string)
	if !isString {
		return controller.NotAuthorized()
	}

	this.Set("authorization", authorization)
	this.Set("access", this.Session().Clear("access"))
	this.Set("refresh", this.Session().Clear("refresh"))
	this.Set("expiry", this.Session().Clear("expiry"))
	this.Set("auth_id", this.Session().Clear("auth_id"))

	this.Set("Title", "New User")

	return this.DefaultResponse()
}

func (this UserController) Create() (response controller.Response) {
	var authData user.AuthorizationData
	var err error

	//would be nice to replace this with some kind of reflection based reader
	authData.Authorization = this.GetFormValue("User[Authorization][Type]")
	authData.Id = this.GetFormValue("User[Authorization][ID]")
	authData.Token.AccessToken = this.GetFormValue("User[Authorization][Access]")
	authData.Token.RefreshToken = this.GetFormValue("User[Authorization][Refresh]")
	authData.Token.Expiry, err = time.Parse(time.RFC1123, this.GetFormValue("User[Authorization][Expiry]"))
	if err != nil {
		panic("Can't Convert Expiry to Time")
	}

	defer func() {
		rec := recover()
		if rec != nil {
			response = controller.FromRack(login.NewUserForm(authData).Run(this.GetRackFuncVars()))
		}
	}()

	var u user.User = user.NewUser()

	u.ClashTag = this.GetFormValue("User[ClashTag]")

	u.Authorizations = []user.AuthorizationData{authData}

	errs := model.Save(u)
	if errs != nil {
		panic(errs)
	}

	this.Apply(login.LogIn(u))
	return this.RespondWith(u)
}

func init() {
	controller.RegisterController(&UserController{u: user.U}).AddToRoot()
}
