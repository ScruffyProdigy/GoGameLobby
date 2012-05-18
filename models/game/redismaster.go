package game

import (
	"../../gamedata"
	"errors"
	"fmt"
)

func startMode(game, mode string, m Mode) (restart bool) {
	var startInfo gamedata.StartInfo
	var err error
	start := QueueMutex.Write.Try(func() {
		preStartInfo := make(map[string]map[string]string)
		for group, full := range m.GroupCount {
			fmt.Print("-Forming", group)
			preStartInfo[group] = make(map[string]string)
			for i := 0; i < full; i++ {
				user := pullFromQueue(game, mode, group)
				join := takeJoinOptions(user, game, mode)
				fmt.Println("--Adding", user, "with", join)
				preStartInfo[group][user] = join
				removeFromMode(user, game, mode)
			}
		}

		model := G.GameFromName(game)
		startInfo, err = model.StartClash(mode, preStartInfo)
		if err != nil {
			fmt.Println("Error:", err.Error())
			//there was an error trying to start the clash, put everybody back in the queue
			for group, players := range preStartInfo {
				for user, join := range players {
					addToMode(user, game, mode)
					jumpTheQueue(user, game, mode, group)
					setJoinOptions(user, game, mode, join)
				}
			}
		} else {
			fmt.Println("Success!")
		}
	})

	if err != nil {
		return
	}

	if start {
		// we've loaded the players in the queue, time to start it
		sendStartMessages(startInfo)
		return false
	}

	// somebody else was trying to start a game, loop through again to check to see if the game can start
	return true
}

func checkForStart(game, mode string, m Mode) {

	start := true
	trying := true

	for trying {
		//if the queue is long enough to start, start it
		QueueMutex.Read.Force(func() {
			fmt.Println("Checking Start for", game, mode)
			for group, full := range m.GroupCount {
				count := int(queueLength(game, mode, group))
				fmt.Println("Testing group:", group)
				fmt.Println(count, "vs", full)
				if count < full {
					start = false
					return
				}
			}
		})

		if !start {
			return
		}

		fmt.Println("starting!")
		trying = startMode(game, mode, m)
	}

}

func AddToQueue(user, game, mode, group, join string, m Mode) (err error) {

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
	isNew := addToMode(user, game, mode)
	if !isNew {
		//the player is already part of the game, remove them from their old queue first
		for group, _ := range m.GroupCount {
			if removeFromQueue(user, game, mode, group) {
				break
			}
		}
	}

	//add the user to the queue
	addToQueue(user, game, mode, group)
	setJoinOptions(user, game, mode, join)

	go checkForStart(game, mode, m)

	return nil
}

func RemoveFromQueue(user, game, mode string, m Mode) (err error) {
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

	if !removeFromMode(user, game, mode) {
		//not in this mode
		return
	}
	for group, _ := range m.GroupCount {
		if removeFromQueue(user, game, mode, group) {
			break
		}
	}
	takeJoinOptions(user, game, mode)

	return nil
}
