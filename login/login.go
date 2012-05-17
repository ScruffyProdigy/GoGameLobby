package login

import (
	"../models/user"
	"github.com/HairyMezican/Middleware/oauther"
	"github.com/HairyMezican/Middleware/redirecter"
	"github.com/HairyMezican/Middleware/sessioner"
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheRack/rack"
	"github.com/HairyMezican/goauth2/oauth"
	"net/http"
	"time"
)

type NewUserForm user.AuthorizationData

func (this NewUserForm) Run(vars map[string]interface{}, next func()) {
	s := sessioner.V(vars)
	s.Set("authorization", this.Authorization)
	s.Set("auth_id", this.Id)
	s.Set("access", this.Token.AccessToken)
	s.Set("refresh", this.Token.RefreshToken)
	s.Set("expiry", this.Token.Expiry.Format(time.RFC1123))

	(redirecter.V)(vars).Redirect("/users/new")
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

func HandleToken(o oauther.Oauther, tok *oauth.Token) rack.Middleware {
	this := new(TokenHandler)
	this.a = o.(Authorizer)
	this.tok = tok
	return this
}

func (this TokenHandler) Run(vars map[string]interface{}, next func()) {
	if this.tok == nil {
		sessioner.V(vars).AddFlash("Log In Cancelled")
		(redirecter.V)(vars).Redirect("/")
		return
	}

	authorization := this.a.GetName()
	id := this.a.GetUserID(this.tok)

	u := user.U.UserFromAuthorization(authorization, id)
	if u != nil {
		//if we find a user, log them in
		(V)(vars).LogIn(u)
		(sessioner.V)(vars).AddFlash("Welcome back, " + u.ClashTag)
		(redirecter.V)(vars).Redirect("/")
	} else {
		//otherwise, have them fill out the New User Form!
		NewUserForm{Token: *this.tok, Authorization: authorization, Id: id}.Run(vars, next)
	}
}

type V map[string]interface{}

func (vars V) CurrentUser() (u *user.User, loggedIn bool) {
	u, loggedIn = vars["CurrentUser"].(*user.User)
	return
}

func (vars V) LogIn(u *user.User) {
	if u == nil {
		vars.LogOut()
	} else {
		sessioner.V(vars).Set("CurrentUser", u.ClashTag)
		vars["CurrentUser"] = u
	}
}

func (vars V) LogInFromClashTag(clashtag string) {
	u := user.U.UserFromClashTag(clashtag)
	vars.LogIn(u)
}

func (vars V) LogOut() {
	(sessioner.V)(vars).Clear("CurrentUser")
	delete(vars, "CurrentUser")
}

var LogOut rack.Func = func(vars map[string]interface{}, next func()) {
	(V)(vars).LogOut()
	(sessioner.V)(vars).AddFlash("You Have Now Logged Out")
	w := (httper.V)(vars).BlankResponse()
	r := (httper.V)(vars).GetRequest()
	http.Redirect(w, r, "/", http.StatusFound)
	w.Save()
}

var Middleware rack.Func = func(vars map[string]interface{}, next func()) {
	u := (sessioner.V)(vars).Get("CurrentUser")
	if u != nil {
		clashtag, isString := u.(string)
		if isString {
			V(vars).LogInFromClashTag(clashtag)
		}
	}
	next()
}
