/*
	Oather provides an interface for any OAuth service to work in with a rack based system
*/

package oauther

import (
	"../goauth2/oauth"
	"../interceptor"
	"../log"
	"../rack"
	"../session"
	"net/http"
)

type TokenHandler interface {
	HandleToken(o Oauther, tok *oauth.Token) rack.Middleware
}

type Oauther interface {
	GetConfig() *oauth.Config
	GetStartUrl() string
	GetRedirectUrl() string
}

func GetCode(o Oauther) rack.Middleware {
	return func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
		state := "" //should be random and arbitrary
		vars.Apply(session.Set("state", state))
		w := rack.BlankResponse()
		url := o.GetConfig().AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusFound)
		return w.Results()
	}
}

func GetToken(o Oauther, U TokenHandler) rack.Middleware {
	return func(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
		//Step 1: Ensure states match
		state1 := r.FormValue("state")
		state2 := vars.Apply(session.Clear("state"))
		if state1 != state2 {
			//states don't match; potential CSRF attempt, we're just going to pass it on, and a 404 will probably be passed back (unless this happens to route somewhere else too)
			//perhaps we should just return a 401-Unauthorized, though
			log.Warning("Potential CSRF attempt")
			return next()
		}

		//Step 2: Exchange the code for the token
		code := r.FormValue("code")
		t := &oauth.Transport{oauth.Config: o.GetConfig()}
		tok, _ := t.Exchange(code)

		//Step 3: Have some other middleware handle whatever they're doing with the token (probably logging a user in)
		process := U.HandleToken(o, tok)
		return process(r, vars, next)
	}
}

func RegisterOauth(i interceptor.Interceptor, o Oauther, U TokenHandler) {
	i.Intercept(o.GetStartUrl(), GetCode(o))
	i.Intercept(o.GetRedirectUrl(), GetToken(o, U))
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
