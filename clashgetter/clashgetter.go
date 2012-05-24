package clashgetter

import (
	"../login"
	"../models/game"
)

type Middleware struct {
}

func (Middleware) Run(vars map[string]interface{}, next func()) {
	currentUser, loggedIn := (login.V)(vars).CurrentUser()
	if loggedIn {
		queues := game.GetUserQueues(currentUser.ClashTag)
		vars["Queues"] = queues

		clashes := game.GetUserClashes(currentUser.ClashTag)
		vars["Clashes"] = clashes
	}
	next()
}

var QueueGetter Middleware
