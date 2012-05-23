package game

import (
	"../../gamedata"
	"../../pubsuber"
	"../../redis"
	"../../websocketcontrol"
	"../user"
	"fmt"
	neturl "net/url"
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

		url := getActiveGameUrl(u.ClashTag, game, mode)
		result := map[string]StartLoc{"startloc": StartLoc{Loc: url}}
		fmt.Println(result)
		return result
	})
}

func usersIndex(user string) string {
	return "users" + sEp + user + sEp + "queues"
}

func queueIndex(game, mode, group string) string {
	return "queues" + sEp + game + sEp + mode + sEp + "groups" + sEp + group
}

func joinIndex(game, mode, user string) string {
	return "queues" + sEp + game + sEp + mode + sEp + "players" + sEp + user
}

func activeGameIndex(user, game, mode string) string {
	return "activegames" + sEp + user + sEp + game + sEp + mode
}

func sendStartMessage(user string, game, mode, url string, options map[string]string) {
	urlvals := make(neturl.Values)
	for k, v := range options {
		urlvals.Add(k, v)
	}

	query := urlvals.Encode()
	if query != "" {
		url += "?" + query
	}

	addActiveGame(user, game, mode, url)
	pubsuber.User(user).SendMessage(map[string]Start{"start": Start{Game: game, Mode: mode}})
}

func sendStartMessages(startInfo gamedata.StartInfo, game, mode string) {
	for p, vals := range startInfo.Players {
		go sendStartMessage(p, game, mode, startInfo.Url, vals.UrlValues)
	}
}

func addActiveGame(user, game, mode, url string) {
	original, err := redis.Client.Getset(activeGameIndex(user, game, mode), url)
	if err != nil {
		panic(err)
	}
	if original.String() != "" {
		//		panic("overwriting active game")
	}
}

func getActiveGameUrl(user, game, mode string) string {
	url, err := redis.Client.Get(activeGameIndex(user, game, mode))
	if err != nil {
		panic(err)
	}
	return url.String()
}

func removeActiveGame(user, game, mode string) {
	_, err := redis.Client.Del(activeGameIndex(user, game, mode))
	if err != nil {
		panic(err)
	}
}

type UserQueue struct {
	Game *Game
	Mode string
}

func GetUserQueues(user string) (result []UserQueue) {
	index := usersIndex(user)
	reply, err := redis.Client.Smembers(index)
	if err != nil {
		panic(err)
	}

	queues := reply.StringArray()

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

func addToMode(user, game, mode string) bool {
	isNew, err := redis.Client.Sadd(usersIndex(user), game+sEp+mode)
	if err != nil {
		panic(err)
	}
	return isNew
}

func removeFromMode(user, game, mode string) bool {
	removed, err := redis.Client.Srem(usersIndex(user), game+sEp+mode)
	if err != nil {
		panic(err)
	}
	return removed
}

func removeFromQueue(user, game, mode, group string) bool {
	numRemoved, err := redis.Client.Lrem(queueIndex(game, mode, group), 0, user)
	if err != nil {
		panic(err)
	}
	return numRemoved > 0
}

func addToQueue(user, game, mode, group string) int64 {
	len, err := redis.Client.Lpush(queueIndex(game, mode, group), user)
	if err != nil {
		panic(err)
	}
	return len
}

func jumpTheQueue(user, game, mode, group string) int64 {
	len, err := redis.Client.Rpush(queueIndex(game, mode, group), user)
	if err != nil {
		panic(err)
	}
	return len
}

func pullFromQueue(game, mode, group string) string {
	user, err := redis.Client.Rpop(queueIndex(game, mode, group))
	if err != nil {
		panic(err)
	}
	return user.String()
}

func queueLength(game, mode, group string) int64 {
	len, err := redis.Client.Llen(queueIndex(game, mode, group))
	if err != nil {
		panic(err)
	}
	return len
}

func setJoinOptions(user, game, mode, join string) {
	err := redis.Client.Set(joinIndex(game, mode, user), join)
	if err != nil {
		panic(err)
	}
}

func takeJoinOptions(user, game, mode string) string {
	join, err := redis.Client.Get(joinIndex(game, mode, user))
	if err != nil {
		panic(err)
	}

	_, err = redis.Client.Del(joinIndex(game, mode, user))
	if err != nil {
		panic(err)
	}
	return join.String()
}
