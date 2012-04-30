package login

import (
	"github.com/HairyMezican/Middleware/oauther"
	"github.com/HairyMezican/Middleware/oauther/facebooker"
	"github.com/HairyMezican/goauth2/oauth"
	"encoding/json"
	"net/http"
)

type Facebooker struct {
	*facebooker.Facebooker
}

func (Facebooker) GetName() string {
	return "facebook"
}

func (this Facebooker) GetUserID(tok *oauth.Token) (result string) {
	oauther.GetSite(this, tok, "https://graph.facebook.com/me", func(res *http.Response) {
        //use json to read in the result into this struct
        var uid struct {
            ID string `json:"id"` //there are a lot of fields, but we really only care about the ID
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

func NewFacebooker(data facebooker.Data) Facebooker {
	return Facebooker{facebooker.New(data)}
}