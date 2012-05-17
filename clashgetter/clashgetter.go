package clashgetter

import (
	"../login"
	"../models/game"
)

type Middleware struct {
}

func (Middleware) Run(vars map[string]interface{}, next func()) {
	currentUser, _ := (login.V)(vars).CurrentUser()
	queues := game.GetUserQueues(currentUser.ClashTag)
	vars["Queues"] = queues
	next()
}

var QueueGetter Middleware
