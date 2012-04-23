package controller

import (
	"../log"
	"../login"
	"../models/user"
	"../models"
	"../rack"
	"../redirecter"
	"../routes"
	"../session"
	"net/http"
	"strconv"
	"time"
	"fmt"
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

func (this UserController) Indexer(s string)  (interface{},bool) {
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
	
	fmt.Fprint(log.DebugLog(),"\n Debug - looking at User:",u)

	vars["Title"] = u.ClashTag

	return next()
}

func (this UserController) New(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	vars["Title"] = "New User"

	authorization, isString := vars.Apply(session.Clear("authorization")).(string)
	if !isString {
		log.Warning("No Authorization found")
		return http.StatusUnauthorized,make(http.Header),[]byte("")
	}

	vars["authorization"] = authorization
	vars["access"] = vars.Apply(session.Clear("access"))
	vars["refresh"] = vars.Apply(session.Clear("refresh"))
	vars["expiry"] = vars.Apply(session.Clear("expiry"))
	vars["auth_id"] = vars.Apply(session.Clear("auth_id"))

	return next()
}

func (this UserController) Create(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	var u user.User
	var data user.AuthorizationData
	defer func() {
		rec := recover()
		if rec != nil {
			status, header, message = login.NewUserForm{Authorization: data.Authorization, ID: data.Id, Tok: &data.Token}.Run(r, vars, next)
		}
	}()

	u.ClashTag = r.FormValue("User[ClashTag]")
	u.Points, err = strconv.Atoi(r.FormValue("User[Points]"))
	if err != nil {
		panic(err)
	}

	//would be nice to replace this with some kind of reflection based reader
	data.Authorization = r.FormValue("User[Authorization][Type]")
	data.Id = r.FormValue("User[Authorization][ID]")
	data.Token.AccessToken = r.FormValue("User[Authorization][Access]")
	data.Token.RefreshToken = r.FormValue("User[Authorization][Refresh]")
	data.Token.Expiry, err = time.Parse(time.RFC1123, r.FormValue("User[Authorization][Expiry]"))
	if err != nil {
		panic("Can't Convert Expiry to Time")
	}

	u.Authorizations = []user.AuthorizationData{data}
	
	errs := model.Save(&u)
	if errs != nil {
		panic(errs)
	}

	return redirecter.Go(r,vars,u.Url(), login.LogIn(&u))
}

func init() {
	routes.Resource(UserController{user.U}).AddTo(routes.Root)
}
