package game

import (
	"../../websocketcontrol"
	"../clash"
	"../user"
	"strings"
)

type Start struct {
	Game string `json:"game"`
	Mode string `json:"mode"`
}

type StartLoc struct {
	Loc string `json:"loc"`
}

func init() {
	websocketcontrol.AddLogoutChore(func(username string) {
		RemoveFromAllQueues(username)
	})
}

type UserClash struct {
	Game string
	Mode string
	Url  string
}

func GetUserClashes(user string) (result []UserClash) {
	clashes := clash.FromUser(user)

	result = make([]UserClash, 0, len(clashes))
	for _, c := range clashes {
		game, mode, url := c.Details(user)

		clashStruct := UserClash{
			Game: game,
			Mode: mode,
			Url:  url,
		}
		result = append(result, clashStruct)
	}
	return
}

type UserQueue struct {
	Game *Game
	Mode string
}

func GetUserQueues(u string) (result []UserQueue) {
	queues := <-(&user.User{ClashTag: u}).Queues().Members()

	result = make([]UserQueue, 0, len(queues))
	for _, queue := range queues {
		queueInfo := strings.SplitN(queue, sEp, 2)
		queueStruct := UserQueue{
			Game: G.GameFromName(queueInfo[0]),
			Mode: queueInfo[1],
		}
		result = append(result, queueStruct)
	}
	return
}
