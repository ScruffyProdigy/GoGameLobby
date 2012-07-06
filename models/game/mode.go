package game

import (
	"../../gamedata"
	"errors"
	"fmt"
	"github.com/HairyMezican/SimpleRedis/redis"
	"../user"
	"../clash"
	"../../global"
)

type Mode struct {
	game       string
	mode       string
	GroupCount map[string]int //a list of the groups needed for the mode, and the number of people needed to fill the group
}

func (m Mode) GameMode() string {
	return m.game + sEp + m.mode
}

func (m Mode) Mode() string {
	return m.mode
}

func (m Mode) Game() string {
	return m.game
}

func (m Mode) gatherPlayers() map[string]map[string]string {
	players := make(map[string]map[string]string)
	for group, full := range m.GroupCount {
		players[group] = make(map[string]string)
		for i := 0; i < full; i++ {
			u := <-m.queues(group).LeftPop()
			join := <-m.joinData(u).Get()
			m.joinData(u).Delete()
			players[group][u] = join
			(&user.User{ClashTag:u}).Queues().Remove(m.GameMode())
		}
	}
	return players
}

func (m Mode) ungatherPlayers(players map[string]map[string]string) {
	for group, players := range players {
		for u, join := range players {
			(&user.User{ClashTag:u}).Queues().Add(m.GameMode())
			m.queues(group).LeftPush(u)
			m.joinData(u).Set(join)
		}
	}
}

func (m Mode) start() (restart bool) {
	var startInfo gamedata.StartInfo
	var err error
	start := QueueMutex.Write.Try(func() {
		players := m.gatherPlayers()

		model := G.GameFromName(m.game)
		startInfo, err = model.StartClash(m.mode, players)
		if err != nil {
			//there was an error trying to start the clash, put everybody back in the queue
			m.ungatherPlayers(players)
		}
	})

	if err != nil {
		return
	}

	if start {
		// we've loaded the players in the queue, time to start it
		clash.New(m,startInfo)
		return false
	}

	// somebody else was trying to start a game, loop through again to check to see if the game can start
	return true
}

func (m Mode) checkForStart() {

	start := true
	trying := true

	for trying {
		//if the queue is long enough to start, start it
		QueueMutex.Read.Force(func() {
			for group, full := range m.GroupCount {
				count := <-m.queues(group).Length()
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

func (m Mode) AddToQueue(u, group, join string) (err error) {

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
	isNew := <-(&user.User{ClashTag:u}).Queues().Add(m.GameMode())
	if !isNew {
		//the player is already part of the game, remove them from their old queue first
		for group, _ := range m.GroupCount {
			if <-m.queues(group).Remove(u) > 0 {
				break
			}
		}
	}

	//add the user to the queue
	m.queues(group).RightPush(u)
	m.joinData(u).Set(join)

	go m.checkForStart()

	return nil
}

func (m Mode) RemoveFromQueue(u string) (err error) {
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

	if <-(&user.User{ClashTag:u}).Queues().Remove(m.GameMode()) {
		//not in this mode
		return
	}
	for group, _ := range m.GroupCount {
		if <-m.queues(group).Remove(u) > 0 {
			break
		}
	}
	m.joinData(u).Delete()

	return nil
}

func (m Mode) prefix() redis.Prefix {
	return global.Redis.Prefix("queues "+m.GameMode()+" ")
}

func (m Mode) queues(group string) redis.List {
	return m.prefix().List("groups " + group)
}

func (m Mode) joinData(u string) redis.String {
	return m.prefix().String("players " + u)
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
