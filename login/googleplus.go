package login

import(
	"github.com/HairyMezican/Middleware/oauther"
	"github.com/HairyMezican/Middleware/oauther/googleplusser"
	"github.com/HairyMezican/goauth2/oauth"
	"net/http"
	"encoding/json"
)

type GooglePlusser struct {
	*googleplusser.GooglePlus
	apikey string
}

func (GooglePlusser) GetName() string {
	return "google"
}

func (this GooglePlusser) GetUserID(tok *oauth.Token) (result string) {
	oauther.GetSite(this, tok, "https://www.googleapis.com/plus/v1/people/me?key="+this.apikey, func(res *http.Response) {
		//use json to read in the result, and get 
		var uid struct {
			ID string `json:"id"` //this is really the only field we care about, we don't really care where people work or any of that shit
			//perhaps in the future, we will take in the age field or something, so we can get a better idea of who our demographics are and cater to them better
			//but for now, we don't really give a fuck
		}

		d := json.NewDecoder(res.Body)
		err := d.Decode(&uid)
		if err != nil {
			panic(err)
		}

		result = uid.ID
	})
	return
}

func NewGooglePlusser(data googleplusser.Data) GooglePlusser {
	return GooglePlusser{googleplusser.New(data),data.Apikey}
}