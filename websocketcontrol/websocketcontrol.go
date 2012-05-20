package websocketcontrol

import (
	"../login"
	"../models/game"
	"../pubsuber"
	"../redis"
	"../trigger"
	"github.com/HairyMezican/Middleware/websocketer"
	"github.com/HairyMezican/TheRack/httper"
	"io"
	"time"
)

const (
	closerIndex = "websocket:closer"
)

func cancelLogoutIndex(user string) string {
	return "Cancel Logout " + user
}

func logout(user string) {
	game.RemoveFromAllQueues(user)
}

func startLogoutProcess(user string) {
	//wait 5 seconds, and if we haven't gotten a cancel message, log the user out
	canceler, err := redis.Subscribe(cancelLogoutIndex(user))
	if err != nil {
		//we weren't able to check to see if the user could log out, so we will automatically log him out now
		logout(user)
		return
	}

	//cleanup at the end
	defer canceler.Close()
	defer canceler.Unsubscribe(cancelLogoutIndex(user))

	//if we don't get a cancel message before we get the 5second callback, log the user out
	select {
	case <-canceler.Messages:
		return
	case <-time.After(5 * time.Second):
		logout(user)
	}
}

type Opener struct{}

func (Opener) Run(vars map[string]interface{}, next func()) {
	currentUser, loggedIn := (login.V)(vars).CurrentUser()
	r := (httper.V)(vars).GetRequest()

	t := trigger.New()

	sendBasicMessage := func(message string) {
		(websocketer.V)(vars).SendBasicMessage(message)
	}

	if loggedIn {
		redis.Client.Publish(cancelLogoutIndex(currentUser.ClashTag), "Doesn't matter what goes here")
		closer := pubsuber.User(currentUser.ClashTag).ReceiveMessages(sendBasicMessage)
		t.OnClose(func() {
			closer.Close()
		})
	}
	closer := pubsuber.Url(r.URL.String()).ReceiveMessages(sendBasicMessage)
	t.OnClose(func() {
		closer.Close()
	})

	vars[closerIndex] = t
}

type Closer struct{}

func (Closer) Run(vars map[string]interface{}, next func()) {
	closer := vars[closerIndex].(io.Closer)
	closer.Close()

	currentUser, loggedIn := (login.V)(vars).CurrentUser()
	if loggedIn {
		if !pubsuber.User(currentUser.ClashTag).IsActive() {
			go startLogoutProcess(currentUser.ClashTag)
		}
	}
}

var OpenUp Opener
var CloseDown Closer
