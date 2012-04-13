/*
	Oather provides an interface for any OAuth service to work in with a rack based system
*/

package oauther

import (
	"../goauth2/oauth"
	"../interceptor"
	"../rack"
	"../session"
	"net/http"
)

type Oauther interface {
	GetConfig() *oauth.Config
	GetAuthUrl() string
	GetRedirectUrl() string
	PreTokenCallback(r *http.Request, vars rack.Vars)
	PostTokenCallback(r *http.Request, vars rack.Vars, tok *oauth.Token) (int, http.Header, []byte)
}

func GetCode(o Oauther) rack.Middleware {
	return func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
		o.PreTokenCallback(r, vars)
		state := "" //should be random and arbitrary
		vars.Apply(session.Set("state", state))
		w := rack.BlankResponse()
		url := o.GetConfig().AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusFound)
		return w.Results()
	}
}

func GetToken(o Oauther) rack.Middleware {
	return func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
		state1 := r.FormValue("state")
		state2 := vars.Apply(session.Clear("state"))
		if state1 != state2 {
			panic("states don't match!  Potential CSRF issue!")
		}

		code := r.FormValue("code")
		t := &oauth.Transport{oauth.Config: o.GetConfig()}
		tok, err := t.Exchange(code)
		if err != nil {
			panic(err)
		}
		return o.PostTokenCallback(r, vars, tok)
	}
}

func RegisterOauth(i interceptor.Interceptor, o Oauther) {
	i.Intercept(o.GetAuthUrl(), GetCode(o))
	i.Intercept(o.GetRedirectUrl(), GetToken(o))
}

func GetSite(o Oauther, tok *oauth.Token, site string, handler func(*http.Response)) {
	t := &oauth.Transport{oauth.Config: o.GetConfig(), oauth.Token: tok}
	req, err := t.Client().Get(site)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	handler(req)
}
