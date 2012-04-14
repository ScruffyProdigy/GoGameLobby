package login

import (
	"../interceptor"
	"../models/user"
	"../rack"
	"../session"
	"net/http"
)

func CurrentUser() rack.VarFunc {
	return session.Get("CurrentUser")
}

func LogIn(u *user.User) rack.VarFunc {
	return func(vars rack.Vars) interface{} {
		vars.Apply(session.Set("CurrentUser", u.ClashTag))
		vars["CurrentUser"] = u
		return nil
	}
}

func LogInFromClashTag(clashtag string) rack.VarFunc {
	u := user.U.UserFromClashTag(clashtag)
	return LogIn(u)
}

func LogOut() rack.VarFunc {
	return func(vars rack.Vars) interface{} {
		vars.Apply(session.Clear("CurrentUser"))
		delete(vars, "CurrentUser")
		return nil
	}
}

func RegisterLogout(i interceptor.Interceptor, url string) {
	i.Intercept(url, func(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
		vars.Apply(LogOut())
		vars.Apply(session.AddFlash("You Have Now Logged Out"))
		w := rack.BlankResponse()
		http.Redirect(w, r, "/", http.StatusFound)
		return w.Results()
	})
}

func Middleware(r *http.Request, vars rack.Vars, next rack.NextFunc) (int, http.Header, []byte) {
	u := vars.Apply(session.Get("CurrentUser"))
	if u != nil {
		clashtag, isString := u.(string)
		if isString {
			vars.Apply(LogInFromClashTag(clashtag))
		}
	}
	return next()
}
