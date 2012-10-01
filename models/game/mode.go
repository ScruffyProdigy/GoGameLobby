package game

import (
	"../../global"
	"../user"
	"errors"
	"fmt"
	"github.com/HairyMezican/SimpleRedis/redis"
)

type Mode struct {
	parent     *Game
	Game       string
	Mode       string
	GroupCount map[string]int //a list of the groups needed for the mode, and the number of people needed to fill the group
}

//New Mode -> Game.newMode

func (m Mode) GameMode() string {
	return m.Game + sEp + m.Mode
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
			(&user.User{ClashTag: u}).Queues().Remove(m.GameMode())
		}
	}
	return players
}

func (m Mode) ungatherPlayers(players map[string]map[string]string) {
	for group, players := range players {
		for u, join := range players {
			(&user.User{ClashTag: u}).Queues().Add(m.GameMode())
			m.queues(group).LeftPush(u)
			m.joinData(u).Set(join)
		}
	}
}

func (m Mode) checkForStart() {
	isFull := true

	for { //keep trying until it works, or we find out that it can't work
		//if the queue is long enough to start, start it
		m.queueMutex().Read.Force(func() {
			for group, full := range m.GroupCount {
				count := <-m.queues(group).Length()
				fmt.Println("Mode:", m.Mode, count, "/", full)
				if count < full {
					isFull = false
				}
			}
		})
		if isFull {
			if m.queueMutex().Write.Try(func() { //this could be the source of some problems when the server starts getting busy
				m.startClash()
			}) {
				//succeeded in trying to make a clash
				return
			} else {
			}
		} else {
			return
		}
	}
}

func (m Mode) startClash() {
	players := m.gatherPlayers()
	err := m.parent.StartClash(m.Mode, players)
	if err != nil {
		m.ungatherPlayers(players)
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
	isNew := <-(&user.User{ClashTag: u}).Queues().Add(m.GameMode())
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

	if <-(&user.User{ClashTag: u}).Queues().Remove(m.GameMode()) {
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
	return global.Redis.Prefix("queues " + m.GameMode() + " ")
}

func (m Mode) queues(group string) redis.List {
	return m.prefix().List("groups " + group)
}

func (m Mode) joinData(u string) redis.String {
	return m.prefix().String("players " + u)
}

func (m Mode) queueMutex() *redis.ReadWriteMutex {
	return m.prefix().ReadWriteMutex("QueueMutex", 16)
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
