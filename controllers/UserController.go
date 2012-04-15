package controller

import (
	"../login"
	"../models/user"
	"../rack"
	"../routes"
	"../session"
	"net/http"
	"strconv"
	"time"
)

var U = user.U

func init() {
	rest := map[string]routes.HandlerFunc{
		"index": func(res routes.Responder, req *http.Request, vars rack.Vars) {
			var users []user.User
			err := U.AllUsers(&users)
			if err != nil {
				panic(err)
			}

			vars["Users"] = users
			vars["Title"] = "Users"
			vars["Layout"] = "base"

			res.Render("users/index")
		},
		"show": func(res routes.Responder, req *http.Request, vars rack.Vars) {
			u := vars["User"].(*user.User)
			vars["Title"] = u.ClashTag
			vars["Layout"] = "base"
			res.Render("users/show")
		},
		"new": func(res routes.Responder, req *http.Request, vars rack.Vars) {
			vars["Title"] = "New User"
			vars["Layout"] = "base"

			authorization, isString := vars.Apply(session.Clear("authorization")).(string)
			if !isString {
				vars.Apply(session.AddFlash("Please Log In with one of the specified providers"))
				res.RedirectTo(routes.Url("/"))
				return
			}

			vars["authorization"] = authorization
			vars["access"] = vars.Apply(session.Clear("access"))
			vars["refresh"] = vars.Apply(session.Clear("refresh"))
			vars["expiry"] = vars.Apply(session.Clear("expiry"))
			vars["auth_id"] = vars.Apply(session.Clear("auth_id"))

			res.Render("users/new")
		},
		"create": func(res routes.Responder, req *http.Request, vars rack.Vars) {
			err := req.ParseForm()
			if err != nil {
				panic(err)
			}

			var u user.User
			u.ClashTag = req.FormValue("User[ClashTag]")
			u.Points, err = strconv.Atoi(req.FormValue("User[Points]"))
			if err != nil {
				panic(err)
			}

			var data user.AuthorizationData
			data.Authorization = req.FormValue("User[Authorization][Type]")
			data.Id = req.FormValue("User[Authorization][ID]")
			data.Token.AccessToken = req.FormValue("User[Authorization][Access]")
			data.Token.RefreshToken = req.FormValue("User[Authorization][Refresh]")
			data.Token.Expiry, err = time.Parse(time.RFC1123, req.FormValue("User[Authorization][Expiry]"))
			if err != nil {
				panic("Can't Convert Expiry to Time")
			}
			u.Authorizations = []user.AuthorizationData{data}

			err = U.AddUser(&u)
			if err != nil {
				panic(err)
			}

			vars.Apply(login.LogIn(&u))

			res.RedirectTo(u)
		},
	}

	userResource := routes.Resource(U, rest)

	routes.Root.AddRoute(userResource)
}
