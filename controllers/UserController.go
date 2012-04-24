package controller

import (
	"../log"
	"../login"
	"../models"
	"../models/user"
	"../rack"
	"../routes"
	"../session"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type UserController struct {
	U *user.UserCollection
}

func (UserController) RouteName() string {
	return "users"
}

func (UserController) VarName() string {
	return "User"
}

func (this UserController) Indexer(s string, vars rack.Vars) (interface{}, bool) {
	result := this.U.UserFromClashTag(s)
	return result, result != nil
}

func (this UserController) Index(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	var users []user.User
	err := this.U.AllUsers(&users)
	if err != nil {
		panic(err)
	}

	vars["Users"] = users
	vars["Title"] = "Users"

	return next()
}

func (this UserController) Show(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	u := vars["User"].(*user.User)

	fmt.Fprint(log.DebugLog(), "\n Debug - looking at User:", u)

	vars["Title"] = u.ClashTag

	return next()
}

func (this UserController) New(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	vars["Title"] = "New User"

	authorization, isString := vars.Apply(session.Clear("authorization")).(string)
	if !isString {
		log.Warning("No Authorization found")
		return http.StatusUnauthorized, make(http.Header), []byte("")
	}

	vars["authorization"] = authorization
	vars["access"] = vars.Apply(session.Clear("access"))
	vars["refresh"] = vars.Apply(session.Clear("refresh"))
	vars["expiry"] = vars.Apply(session.Clear("expiry"))
	vars["auth_id"] = vars.Apply(session.Clear("auth_id"))

	return next()
}

func (this UserController) Create(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	var authData user.AuthorizationData
	var err error

	//would be nice to replace this with some kind of reflection based reader
	authData.Authorization = r.FormValue("User[Authorization][Type]")
	authData.Id = r.FormValue("User[Authorization][ID]")
	authData.Token.AccessToken = r.FormValue("User[Authorization][Access]")
	authData.Token.RefreshToken = r.FormValue("User[Authorization][Refresh]")
	authData.Token.Expiry, err = time.Parse(time.RFC1123, r.FormValue("User[Authorization][Expiry]"))
	if err != nil {
		panic("Can't Convert Expiry to Time")
	}

	defer func() {
		rec := recover()
		if rec != nil {
			status, header, message = login.NewUserForm{Authorization: authData.Authorization, ID: authData.Id, Tok: &authData.Token}.Run(r, vars, next)
		}
	}()

	var u user.User

	u.ClashTag = r.FormValue("User[ClashTag]")
	u.Points, err = strconv.Atoi(r.FormValue("User[Points]"))
	if err != nil {
		panic(err)
	}

	u.Authorizations = []user.AuthorizationData{authData}

	errs := model.Save(&u)
	if errs != nil {
		panic(errs)
	}

	vars["User"] = &u
	vars.Apply(login.LogIn(&u))
	return next()
}

func init() {
	routes.Resource(UserController{user.U}).AddTo(routes.Root)
}
