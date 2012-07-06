package clashgetter

import (
	"../login"
	"../models/game"
	"../models/clash"
)

type Middleware struct {
}

func (Middleware) Run(vars map[string]interface{}, next func()) {
	currentUser, loggedIn := (login.V)(vars).CurrentUser()
	if loggedIn {
		queues := game.GetUserQueues(currentUser.ClashTag)
		vars["Queues"] = queues

		clashes := clash.FromUser(currentUser.ClashTag)
		vars["Clashes"] = clashes
	}
	next()
}

var QueueGetter Middleware
