/*
	the Facebooker package implements the Oauther interface, and provides facebook specific interactivity
*/
package facebooker

import (
	"../goauth2/oauth"
	"../log"
	"../login"
	"../oauther"
	"encoding/json"
	"net/http"
	"strings"
)

const ()

type Data struct {
	// most of these variables should be straight up copied from https://developers.facebook.com/apps
	// if you mess up any of these variables, you will get an error
	AppId       string   `json:"app_id"`       // on the dashboard - "App ID/API Key"
	AppSecret   string   `json:"app_secret"`   // on the dashboard - "App Secret"
	SiteUrl     string   `json:"site_url"`     // on the dashboard - "Site URL"
	Permissions []string `json:"permissions"`  // see http://developers.facebook.com/docs/authentication/permissions/ for more details
	StartUrl    string   `json:"start_url"`    // you decide - the route where the user should get start
	RedirectUrl string   `json:"redirect_url"` // you decide - where facebook sends the user after they've been authenticated
}

type facebooker struct {
	data   Data
	config *oauth.Config
}

//most applications only need one set of application settings - this is where that should be stored
var Default oauther.Oauther

func NewFacebooker(data Data) login.Authorizer {
	this := new(facebooker)
	this.data = data
	return this
}

func SetConfiguration(data Data) {
	Default = NewFacebooker(data)
}

func (*facebooker) GetName() string {
	return "facebook"
}

func (this *facebooker) GetConfig() *oauth.Config {
	if this.config == nil {
		this.config = new(oauth.Config)
		this.config.ClientId = this.data.AppId
		this.config.ClientSecret = this.data.AppSecret
		this.config.Scope = strings.Join(this.data.Permissions, ",")
		this.config.AuthURL = "https://www.facebook.com/dialog/oauth"
		this.config.TokenURL = "https://graph.facebook.com/oauth/access_token"
		this.config.RedirectURL = this.data.SiteUrl + this.data.RedirectUrl
	}
	return this.config
}

func (this *facebooker) GetStartUrl() string {
	return "/" + this.data.StartUrl
}

func (this *facebooker) GetRedirectUrl() string {
	return "/" + this.data.RedirectUrl
}

func (this *facebooker) GetUserID(tok *oauth.Token) (result string) {
	oauther.GetSite(this, tok, "https://graph.facebook.com/me", func(res *http.Response) {
		//use json to read in the result, and get 
		var uid struct {
			ID string `json:"id"` //this is really the only field we care about, we don't really care where people work or any of that shit
			//perhaps in the future, we will take in the age field or something, so we can get a better idea of who our demographics are and cater to them better
			//but for now, we don't really give a fuck
		}

		d := json.NewDecoder(res.Body)
		err := d.Decode(&uid)
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}

		result = uid.ID
	})
	return
}

func (this *facebooker) GetUserFriends(tok *oauth.Token) (result []string) {
	oauther.GetSite(this, tok, "https://graph.facebook.com/me/friends", func(res *http.Response) {
		//use json to read in the result, and get 
		result = []string{}
	})
	return
}
