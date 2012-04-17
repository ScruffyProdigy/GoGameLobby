package googleplusser

import (
	"../goauth2/oauth"
	"../log"
	"../login"
	"../oauther"
	"encoding/json"
	"net/http"
	"strings"
)

const (
	UserPermission  = "https://www.googleapis.com/auth/plus.me"
	EmailPermission = "https://www.googleapis.com/auth/userinfo.email"
)

type Data struct {
	// most of these variables should be straight up copied from https://code.google.com/apis/console/
	// if you mess up any of these variables, you will get an error
	ClientID     string `json:"client_id"`     // found in OAuth Client ID for Web Applications "Client ID"
	ClientSecret string `json:"client_secret"` // found in OAuth Client ID for Web Applications "Client secret"
	SiteUrl      string `json:"site_url"`      // found in OAuth Client ID for Web Applications "Redirect URIs" (the host of the specified URI)
	RedirectUri  string `json:"redirect_uri"`  // found in OAuth Client ID for Web Applications "Redirect URIs" (the path of the specified URI)
	Apikey       string `json:"api_key"`       // found in Simple API Access (Server Key) "API key"
	// these variables are yours to decide
	StartUri string `json:"start_uri"` // yours to decide, it is the path that you should direct the user to to log in
	//	these variables are created by us
	Permissions []string `json:"permissions"` // what you want to do, options found above "UserPermission is recommended"

}

type googleplus struct {
	data   Data
	config *oauth.Config
}

func NewGooglePlusser(data Data) login.Authorizer {
	gp := new(googleplus)
	gp.data = data
	return gp
}

var Default login.Authorizer

func SetConfiguration(data Data) login.Authorizer {
	Default = NewGooglePlusser(data)
	return Default
}

func (*googleplus) GetName() string {
	return "google"
}

func (this *googleplus) GetStartUrl() string {
	return "/" + this.data.StartUri
}

func (this *googleplus) GetRedirectUrl() string {
	return "/" + this.data.RedirectUri
}

func (this *googleplus) GetConfig() *oauth.Config {
	if this.config == nil {
		this.config = new(oauth.Config)
		this.config.ClientId = this.data.ClientID
		this.config.ClientSecret = this.data.ClientSecret
		this.config.Scope = strings.Join(this.data.Permissions, ",")
		this.config.AuthURL = "https://accounts.google.com/o/oauth2/auth"
		this.config.TokenURL = "https://accounts.google.com/o/oauth2/token"
		this.config.RedirectURL = this.data.SiteUrl + this.data.RedirectUri

	}

	return this.config
}

func (this *googleplus) GetUserID(tok *oauth.Token) (result string) {
	oauther.GetSite(this, tok, "https://www.googleapis.com/plus/v1/people/me?key="+this.data.Apikey, func(res *http.Response) {
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

/*
	The Google+ API currently has no method for figuring out which users somebody is following
	and thus, no way for me to help somebody find their friends on my site :-(

func (this *googleplus) GetUserFriends(tok *oauth.Token) []string {
	//this function is not possible
}	
*/
