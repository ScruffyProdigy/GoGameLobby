package game

import (
	"../../gamedata"
	"../../pubsuber"
	"errors"
	"fmt"
	neturl "net/url"
)

type Mode struct {
	game       string
	mode       string
	GroupCount map[string]int //a list of the groups needed for the mode, and the number of people needed to fill the group
}

func (m Mode) gamemode() string {
	return m.game + sEp + m.mode
}

func (m Mode) start() (restart bool) {
	var startInfo gamedata.StartInfo
	var err error
	start := QueueMutex.Write.Try(func() {
		preStartInfo := make(map[string]map[string]string)
		for group, full := range m.GroupCount {
			preStartInfo[group] = make(map[string]string)
			for i := 0; i < full; i++ {
				user := queues(m.gamemode(), group).LeftPop()
				join, _ := joinData(user, m.gamemode()).Clear()
				preStartInfo[group][user] = join
				userQueues(user).Remove(m.gamemode())
			}
		}

		model := G.GameFromName(m.game)
		startInfo, err = model.StartClash(m.mode, preStartInfo)
		if err != nil {
			//there was an error trying to start the clash, put everybody back in the queue
			for group, players := range preStartInfo {
				for user, join := range players {
					userQueues(user).Add(m.gamemode())
					queues(m.gamemode(), group).LeftPush(user)
					joinData(user, m.gamemode()).Set(join)
				}
			}
		}
	})

	if err != nil {
		return
	}

	if start {
		// we've loaded the players in the queue, time to start it
		m.startClash(startInfo)
		return false
	}

	// somebody else was trying to start a game, loop through again to check to see if the game can start
	return true
}

func (m Mode) addPlayerToClash(user string, url string, options map[string]string) {
	urlvals := make(neturl.Values)
	for k, v := range options {
		urlvals.Add(k, v)
	}

	query := urlvals.Encode()
	if query != "" {
		url += "?" + query
	}

	userClashes(user).Add(m.gamemode())
	clashUrl(user, m.gamemode()).Set(url)
	pubsuber.User(user).SendMessage(map[string]Start{"start": Start{Game: m.game, Mode: m.mode}})
}

func (m Mode) startClash(startInfo gamedata.StartInfo) {
	for p, vals := range startInfo.Players {
		go m.addPlayerToClash(p, startInfo.Url, vals.UrlValues)
	}
}

func (m Mode) checkForStart() {

	start := true
	trying := true

	for trying {
		//if the queue is long enough to start, start it
		QueueMutex.Read.Force(func() {
			for group, full := range m.GroupCount {
				count := int(queues(m.gamemode(), group).Length())
				if count < full {
					start = false
					return
				}
			}
		})

		if !start {
			return
		}

		fmt.Println("starting", m.mode)
		trying = m.start()
	}

}

func (m Mode) AddToQueue(user, group, join string) (err error) {

	defer func() {
		rec := recover()
		if rec != nil {
			if e, ok := rec.(error); ok {
				err = e
				return
			}
			if str, ok := rec.(string); ok {
				err = errors.New(str)
				return
			}
			err = errors.New("UnknownError")
			return
		}
	}()

	//add the user to the set of players queuing for a clash (and find out if they're already on the list)
	isNew := userQueues(user).Add(m.gamemode())
	if !isNew {
		//the player is already part of the game, remove them from their old queue first
		for group, _ := range m.GroupCount {
			if queues(m.gamemode(), group).Remove(user) > 0 {
				break
			}
		}
	}

	//add the user to the queue
	queues(m.gamemode(), group).RightPush(user)
	joinData(user, m.gamemode()).Set(join)

	go m.checkForStart()

	return nil
}

func (m Mode) RemoveFromQueue(user string) (err error) {
	defer func() {
		rec := recover()
		if rec != nil {
			if e, ok := rec.(error); ok {
				err = e
				return
			}
			if str, ok := rec.(string); ok {
				err = errors.New(str)
				return
			}
			err = errors.New("UnknownError")
			return
		}
	}()

	if !userQueues(user).Remove(m.gamemode()) {
		//not in this mode
		return
	}
	for group, _ := range m.GroupCount {
		if queues(m.gamemode(), group).Remove(user) > 0 {
			break
		}
	}
	joinData(user, m.gamemode()).Delete()

	return nil
}

func RemoveFromAllQueues(user string) (err error) {
	fmt.Println("Logging", user, "out from all games")
	queues := GetUserQueues(user)
	for _, queue := range queues {
		mode, _ := queue.Game.GetMode(queue.Mode)
		mode.RemoveFromQueue(user)
	}
	return nil
}
