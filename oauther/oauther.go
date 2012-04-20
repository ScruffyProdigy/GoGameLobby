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
	HandleToken(tok *oauth.Token) rack.Middleware
}

type Oauther interface {
	GetConfig() *oauth.Config
	GetStartUrl() string
	GetRedirectUrl() string
}

type codeGetter struct {
	o Oauther
}

func (this codeGetter) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	state := "" //should be random and arbitrary
	vars.Apply(session.Set("state", state))
	w := rack.BlankResponse()
	url := this.o.GetConfig().AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
	return w.Results()
}

type tokenGetter struct {
	o Oauther
	t TokenHandler
}

func (this tokenGetter) Run(r *http.Request, vars rack.Vars, next rack.NextFunc) (status int, header http.Header, message []byte) {
	//Step 1: Ensure states match
	state1 := r.FormValue("state")
	state2, isString := vars.Apply(session.Clear("state")).(string)

	//if states don't match, it's a potential CSRF attempt; we're just going to pass it on, and a 404 will probably be passed back (unless this happens to route somewhere else too)
	//perhaps we should just return a 401-Unauthorized, though
	if !isString {
		log.Warning("Potential CSRF attempt : cookie not set properly")
		return next()
	}
	if state1 != state2 {
		log.Warning("Potential CSRF attempt : (" + state1 + ") != (" + state2 + ")")
		return next()
	}

	//Step 2: Exchange the code for the token
	code := r.FormValue("code")
	t := &oauth.Transport{oauth.Config: this.o.GetConfig()}
	tok, _ := t.Exchange(code)

	//Step 3: Have some other middleware handle whatever they're doing with the token (probably logging a user in)
	process := this.t.HandleToken(tok)
	return process.Run(r, vars, next)
}

func RegisterOauth(i interceptor.Interceptor, o Oauther, t TokenHandler) {
	i.Intercept(o.GetStartUrl(), &codeGetter{o})
	i.Intercept(o.GetRedirectUrl(), &tokenGetter{o, t})
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
