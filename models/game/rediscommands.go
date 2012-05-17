package game

import (
	"../../gamedata"
	"../../redis"
	neturl "net/url"
	"strings"
)

func usersIndex(user string) string {
	return "users" + sEp + user + sEp + "queues"
}

func queueIndex(game, mode, group string) string {
	return "queues" + sEp + game + sEp + mode + sEp + "groups" + sEp + group
}

func joinIndex(game, mode, user string) string {
	return "queues" + sEp + game + sEp + mode + sEp + "players" + sEp + user
}

func startGameForUser(player string, url string) {
	redis.Client.Publish("Users:"+player, "Start:"+url)
}

func sendStartMessage(user string, url string, options map[string]string) {
	urlvals := make(neturl.Values)
	for k, v := range options {
		urlvals.Add(k, v)
	}

	query := urlvals.Encode()
	if query != "" {
		url += "?" + query
	}

	startGameForUser(user, url)
}

func sendStartMessages(startInfo gamedata.StartInfo) {
	for p, vals := range startInfo.Players {
		go sendStartMessage(p, startInfo.Url, vals.UrlValues)
	}
}

type UserQueue struct {
	Game string
	Mode string
}

func GetUserQueues(user string) (result []UserQueue) {
	reply, err := redis.Client.Smembers(usersIndex(user))
	if err != nil {
		panic(err)
	}

	queues := reply.StringArray()

	result = make([]UserQueue, 0, len(queues))
	for _, queue := range queues {
		queueInfo := strings.SplitN(queue, sEp, 2)
		queueStruct := UserQueue{
			Game: queueInfo[0],
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
