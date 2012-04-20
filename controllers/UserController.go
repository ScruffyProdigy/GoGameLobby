package controller

import (
	"../log"
	"../login"
	"../models/user"
	"../rack"
	"../redirecter"
	"../routes"
	"../session"
	"../templater"
	"net/http"
	"strconv"
	"time"
)

var U = user.U

func init() {
	rest := map[string]rack.Middleware{
		"index": rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
			var users []user.User
			err := U.AllUsers(&users)
			if err != nil {
				panic(err)
			}

			vars["Users"] = users
			vars["Title"] = "Users"
			vars["Layout"] = "base"

			w := rack.BlankResponse()
			templater.Get("users/index").Execute(w, vars)
			return w.Results()
		}),
		"show": rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
			u := vars["User"].(*user.User)

			vars["Title"] = u.ClashTag
			vars["Layout"] = "base"

			lodges := L.LodgesFromMason(u.ClashTag)
			if lodges != nil {
				vars["Lodges"] = lodges
			}

			w := rack.BlankResponse()
			templater.Get("users/show").Execute(w, vars)
			return w.Results()
		}),
		"new": rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
			vars["Title"] = "New User"
			vars["Layout"] = "base"

			authorization, isString := vars.Apply(session.Clear("authorization")).(string)
			if !isString {
				log.Warning("No Authorization found")
				return next()
			}

			vars["authorization"] = authorization
			vars["access"] = vars.Apply(session.Clear("access"))
			vars["refresh"] = vars.Apply(session.Clear("refresh"))
			vars["expiry"] = vars.Apply(session.Clear("expiry"))
			vars["auth_id"] = vars.Apply(session.Clear("auth_id"))

			w := rack.BlankResponse()
			templater.Get("users/new").Execute(w, vars)
			return w.Results()
		}),
		"create": rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
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
			err = U.AddUser(&u)
			if err != nil {
				panic(err)
			}

			return redirecter.Go(u.Url(), login.LogIn(&u)).Run(r, vars, next)
		}),
	}

	userResource := routes.Resource(U, rest)

	routes.Root.AddRoute(userResource.Collection)
}
