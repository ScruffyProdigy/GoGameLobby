/*
	the Facebooker package implements the Oauther interface, and provides facebook specific interactivity
*/
package facebooker

import (
	"../goauth2/oauth"
	"../log"
	"../login"
	"../models/user"
	"../oauther"
	"../rack"
	"../session"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Facebooker interface {
	oauther.Oauther
	GetUserID(*oauth.Token) int
	GetUserFriends(*oauth.Token) []int
}

type facebooker struct {
	// most of these variables should be straight up copied from developer.facebook.com
	// if you mess up any of these variables, you will get an error
	appId       string   // Application ID - supplied by Facebook
	appSecret   string   // Application Secret - supplied by Facebook
	siteUrl     string   // Site URL - you specify this to Facebook
	permissions []string // Permissions - what you want to be able to do on facebook - see http://developers.facebook.com/docs/authentication/permissions/ for more details
	// these variables are yours to decide
	authUrl     string
	redirectUrl string
	//	these variables are created by us
	config *oauth.Config
}

//most applications only need one set of application settings - this is where that should be stored
var Default Facebooker

func NewFacebooker(appId, appSecret, siteUrl, authUrl, redirectUrl string, permissions []string) Facebooker {
	this := new(facebooker)
	this.appId = appId
	this.appSecret = appSecret
	this.siteUrl = siteUrl
	this.permissions = permissions
	this.authUrl = authUrl
	this.redirectUrl = redirectUrl

	return this
}

func SetConfiguration(appId, appSecret, siteUrl, authUrl, redirectUrl string, permissions []string) {
	Default = NewFacebooker(appId, appSecret, siteUrl, authUrl, redirectUrl, permissions)
}

func (this *facebooker) GetConfig() *oauth.Config {
	if this.config != nil {
		return this.config
	}
	this.config = new(oauth.Config)
	this.config.ClientId = this.appId
	this.config.ClientSecret = this.appSecret
	this.config.Scope = strings.Join(this.permissions, ",")
	this.config.AuthURL = "https://www.facebook.com/dialog/oauth"
	this.config.TokenURL = "https://graph.facebook.com/oauth/access_token"
	this.config.RedirectURL = this.siteUrl + this.redirectUrl

	return this.config
}

func (this *facebooker) GetAuthUrl() string {
	return "/" + this.authUrl
}

func (this *facebooker) GetRedirectUrl() string {
	return "/" + this.redirectUrl
}

func (this *facebooker) PreTokenCallback(r *http.Request, vars rack.Vars) {
	// We need to save what was going on, so that afterwards, we can resume
}

func (this *facebooker) PostTokenCallback(r *http.Request, vars rack.Vars, tok *oauth.Token) (status int, header http.Header, message []byte) {
	//save the token in the user
	//and then resume whatever we were doing
	userId := this.GetUserID(tok)
	u := user.U.UserFromFacebookID(userId)
	if u == nil {
		vars.Apply(session.Set("authorization", "facebook"))
		vars.Apply(session.Set("access", tok.AccessToken))
		vars.Apply(session.Set("refresh", tok.RefreshToken))
		vars.Apply(session.Set("expiry", tok.Expiry.Format(time.RFC1123)))
		vars.Apply(session.Set("auth_id", userId))

		w := rack.BlankResponse()
		http.Redirect(w, r, "/users/new", http.StatusFound)
		return w.Results()
	} else {
		vars.Apply(login.LogIn(u))
		vars.Apply(session.AddFlash("Welcome " + u.ClashTag))

		w := rack.BlankResponse()
		http.Redirect(w, r, "/", http.StatusFound)
		return w.Results()
	}
	return
}

type putUserIdHere struct {
	Id string "id"
}

func (this *facebooker) GetUserID(tok *oauth.Token) (result int) {
	oauther.GetSite(this, tok, "https://graph.facebook.com/me", func(res *http.Response) {
		//use json to read in the result, and get 
		var results map[string]interface{}

		d := json.NewDecoder(res.Body)
		err := d.Decode(&results)
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}

		r, isString := results["id"].(string)
		if !isString {
			panic("Facebook has experienced an error, please try logging in again")
		}
		result, err = strconv.Atoi(r)
		if err != nil {
			panic("ID is not a number")
		}
	})
	return
}

func (this *facebooker) GetUserFriends(tok *oauth.Token) (result []int) {
	oauther.GetSite(this, tok, "https://graph.facebook.com/me", func(res *http.Response) {
		//use json to read in the result, and get 
		result = make([]int, 0)
	})
	return
}
