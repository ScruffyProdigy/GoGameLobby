package game

import (
	"github.com/HairyMezican/SimpleRedis/redis"
	"../../global"
	"../../websocketcontrol"
	"strings"
	"../clash"
	"../user"
)

var QueueMutex *redis.ReadWriteMutex

type Start struct {
	Game string `json:"game"`
	Mode string `json:"mode"`
}

type StartLoc struct {
	Loc string `json:"loc"`
}

func init() {
	QueueMutex = global.Redis.ReadWriteMutex("QueueMutex", 16)

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
	for _, c := range(clashes) {
		game,mode,url := c.Details(user)

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
	print("a\n")
	queues := <-(&user.User{ClashTag:u}).Queues().Members()
	print("b\n")

	result = make([]UserQueue, 0, len(queues))
	print("c\n")
	for _, queue := range queues {
		print("d\n")
		queueInfo := strings.SplitN(queue, sEp, 2)
		print("e\n")
		queueStruct := UserQueue{
			Game: G.GameFromName(queueInfo[0]),
			Mode: queueInfo[1],
		}
		print("f\n")
		result = append(result, queueStruct)
		print("g\n")
	}
	print("h\n")
	return
}
