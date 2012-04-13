package login

import (
	"../models/user"
	"../rack"
	"../session"
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
