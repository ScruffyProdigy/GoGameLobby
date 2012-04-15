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
	ClientID     string // Application ID - supplied by Google
	ClientSecret string // Application Secret - supplied by Google
	SiteUrl      string // Site URL - you specify this to Google
	RedirectUri  string //redirect URI - you specify this to Google (this is the part of the Redirect URI that is not also in the javascript origin)
	Apikey       string //API Key - Google requires this along with an OAuth token - in the API console, use a server API key
	// these variables are yours to decide
	StartUri string
	//	these variables are created by us
	Permissions []string // Permissions - Currently there is only one option

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

func (this *googleplus) GetUserFriends(tok *oauth.Token) []string {
	return []string{}
}
