package login

import (
	"../goauth2/oauth"
	"../models/user"
	"../oauther"
	"../rack"
	"../redirecter"
	"../session"
	"net/http"
	"time"
)


type NewUserForm user.AuthorizationData

func (this NewUserForm) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	return redirecter.Go(r, vars, "/users/new",
		session.Set("authorization", this.Authorization),
		session.Set("auth_id", this.Id),
		session.Set("access", this.Token.AccessToken),
		session.Set("refresh", this.Token.RefreshToken),
		session.Set("expiry", this.Token.Expiry.Format(time.RFC1123)),
	)
}

type Authorizer interface {
	oauther.Oauther
	GetName() string
	GetUserID(*oauth.Token) string
}

type TokenHandler struct {
	a   Authorizer
	tok *oauth.Token
}

func CreateHandler(a Authorizer) *TokenHandler {
	t := new(TokenHandler)
	t.a = a
	return t
}

func (this *TokenHandler) HandleToken(tok *oauth.Token) rack.Middleware {
	this.tok = tok
	return this
}

func (this TokenHandler) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	if this.tok == nil {
		vars.Apply(session.AddFlash("Log In Cancelled"))

		w := rack.BlankResponse()
		http.Redirect(w, r, "/", http.StatusFound)
		return w.Results()
	}

	authorization := this.a.GetName()
	id := this.a.GetUserID(this.tok)

	u := user.U.UserFromAuthorization(authorization, id)
	if u != nil {
		//if we find a user, log them in
		vars.Apply(LogIn(u))
		vars.Apply(session.AddFlash("Welcome back, " + u.ClashTag))

		w := rack.BlankResponse()
		http.Redirect(w, r, "/", http.StatusFound)
		return w.Results()
	}

	//otherwise, have them fill out the New User Form!
	return NewUserForm{Token: *this.tok, Authorization: authorization, Id: id}.Run(r, vars, next)
}

func CurrentUser(vars rack.Vars) *user.User {
	return vars["CurrentUser"].(*user.User)
}

func LogIn(u *user.User) rack.VarFunc {
	return func(vars rack.Vars) interface{} {
		if u == nil {
			vars.Apply(LogOutFunc())
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

func LogOutFunc() rack.VarFunc {
	return func(vars rack.Vars) interface{} {
		vars.Apply(session.Clear("CurrentUser"))
		delete(vars, "CurrentUser")
		return nil
	}
}

var LogOut = rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
	vars.Apply(LogOutFunc())
	vars.Apply(session.AddFlash("You Have Now Logged Out"))
	w := rack.BlankResponse()
	http.Redirect(w, r, "/", http.StatusFound)
	return w.Results()
})

var Middleware = rack.Func(func(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
	u := vars.Apply(session.Get("CurrentUser"))
	if u != nil {
		clashtag, isString := u.(string)
		if isString {
			vars.Apply(LogInFromClashTag(clashtag))
		}
	}
	return next()
})
