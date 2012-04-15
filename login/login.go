package login

import (
	"../goauth2/oauth"
	"../interceptor"
	"../models/user"
	"../oauther"
	"../rack"
	"../session"
	"net/http"
	"time"
)

type Authorizer interface {
	oauther.Oauther
	GetName() string
	GetUserID(*oauth.Token) string
	GetUserFriends(*oauth.Token) []string
}

type Middlewarer struct{}

var Logger Middlewarer

func (Middlewarer) HandleToken(o oauther.Oauther, tok *oauth.Token) rack.Middleware {
	return func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
		w := rack.BlankResponse()
		if tok == nil {
			vars.Apply(session.AddFlash("Log In Cancelled"))
			http.Redirect(w, r, "/", http.StatusFound)
			return w.Results()
		}
		a, canAuthorize := o.(Authorizer)
		if !canAuthorize {
			panic("oauth provider doesn't provide user information")
		}

		authorization := a.GetName()
		id := a.GetUserID(tok)

		u := user.U.UserFromAuthorization(authorization, id)
		if u != nil {
			//if we find a user, log them in
			vars.Apply(LogIn(u))
			vars.Apply(session.AddFlash("Welcome back, " + u.ClashTag))

			http.Redirect(w, r, "/", http.StatusFound)
			return w.Results()
		}

		//otherwise, have them fill out the New User Form!
		vars.Apply(session.Set("authorization", authorization))
		vars.Apply(session.Set("access", tok.AccessToken))
		vars.Apply(session.Set("refresh", tok.RefreshToken))
		vars.Apply(session.Set("expiry", tok.Expiry.Format(time.RFC1123)))
		vars.Apply(session.Set("auth_id", id))

		http.Redirect(w, r, "/users/new", http.StatusFound)
		return w.Results()
	}
}

func CurrentUser() rack.VarFunc {
	return session.Get("CurrentUser")
}

func LogIn(u *user.User) rack.VarFunc {
	return func(vars rack.Vars) interface{} {
		if u == nil {
			vars.Apply(LogOut())
			return nil
		}
		vars.Apply(session.Set("CurrentUser", u.ClashTag))
		vars["CurrentUser"] = u
		return nil
	}
}

func LogInFromClashTag(clashtag string) rack.VarFunc {
	u := user.U.UserFromClashTag(clashtag)
	return LogIn(u)
}

func LogOut() rack.VarFunc {
	return func(vars rack.Vars) interface{} {
		vars.Apply(session.Clear("CurrentUser"))
		delete(vars, "CurrentUser")
		return nil
	}
}

func RegisterLogout(i interceptor.Interceptor, url string) {
	i.Intercept(url, func(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
		vars.Apply(LogOut())
		vars.Apply(session.AddFlash("You Have Now Logged Out"))
		w := rack.BlankResponse()
		http.Redirect(w, r, "/", http.StatusFound)
		return w.Results()
	})
}

func Middleware(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
	u := vars.Apply(session.Get("CurrentUser"))
	if u != nil {
		clashtag, isString := u.(string)
		if isString {
			vars.Apply(LogInFromClashTag(clashtag))
		}
	}
	return next()
}
