package game

import (
	"../../redis"
	"../../websocketcontrol"
	"../user"
	"fmt"
	"strings"
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
	var err error
	QueueMutex, err = redis.RWMutex("QueueMutex", 16)
	if err != nil {
		panic("Couldn't set up Queue Mutex")
	}

	websocketcontrol.AddLogoutChore(func(username string) {
		RemoveFromAllQueues(username)
	})

	websocketcontrol.MessageAction("start", func(u *user.User, data interface{}) interface{} {
		startinfo, ok := data.(map[string]interface{})
		if !ok {
			fmt.Println("Couldn't Get Start Info")
			return nil
		}

		game, ok := startinfo["game"].(string)
		if !ok {
			fmt.Println("Couldn't Get Game")
			return nil
		}

		mode, ok := startinfo["mode"].(string)
		if !ok {
			fmt.Println("Couldn't Get Mode")
			return nil
		}

		url := clashUrl(u.ClashTag, game+sEp+mode).Get()
		result := map[string]StartLoc{"startloc": StartLoc{Loc: url}}
		fmt.Println(result)
		return result
	})
}

func userQueues(user string) redis.Set {
	return redis.Set("users" + sEp + user + sEp + "queues")
}

func userClashes(user string) redis.Set {
	return redis.Set("users" + sEp + user + sEp + "clashes")
}

func queues(gamemode, group string) redis.List {
	return redis.List("queues" + sEp + gamemode + sEp + "groups" + sEp + group)
}

func joinData(user, gamemode string) redis.String {
	return redis.String("queues" + sEp + gamemode + sEp + "players" + sEp + user)
}

func clashUrl(user, gamemode string) redis.String {
	return redis.String("clashes" + sEp + user + sEp + gamemode)
}

type UserClash struct {
	Game *Game
	Mode string
	Url  string
}

func GetUserClashes(user string) (result []UserClash) {
	clashes := userClashes(user).Members()

	result = make([]UserClash, 0, len(clashes))
	for _, clash := range clashes {
		clashInfo := strings.SplitN(clash, sEp, 2)

		game := G.GameFromName(clashInfo[0])
		mode := clashInfo[1]
		url := clashUrl(user, clash).Get()

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

func GetUserQueues(user string) (result []UserQueue) {
	queues := userQueues(user).Members()

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
