package websocketcontrol

import (
	"../../Middleware/websocketer"
	"../login"
	"../models/user"
	"../pubsuber"
	"../redis"
	"../trigger"
	"encoding/json"
	"github.com/HairyMezican/Middleware/logger"
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheRack/rack"
	"io"
	"time"
)

const (
	closerIndex = "websocket:closer"
)

func loginChannel(user string) redis.Channel {
	return redis.Redis.Channel("Cancel Logout " + user)
}

var logoutChores = make([]func(string), 0)

func AddLogoutChore(chore func(string)) {
	logoutChores = append(logoutChores, chore)
}

func logout(username string) {
	for _, chore := range logoutChores {
		chore(username)
	}
}

func startLogoutProcess(user string) {
	//wait 5 seconds, and if we haven't gotten a cancel message, log the user out
	loginChannel(user).BlockingSubscription(func(cancellogout <-chan string) {
		//if we don't get a cancel message before we get the 5second callback, log the user out
		select {
		case <-cancellogout:
			return
		case <-time.After(5 * time.Second):
			logout(user)
		}
	})
}

type WebsocketMessage struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func New() rack.Middleware {
	ws := websocketer.New()
	ws.OnOpen(wsOpener)
	ws.OnClose(wsCloser)
	ws.OnMessage(wsMessager)
	ws.UseJSON()
	ws.OnStorage(func() interface{} {
		var m WebsocketMessage
		return &m
	})

	return ws
}

var wsOpener rack.Func = func(vars map[string]interface{}, next func()) {
	currentUser, loggedIn := (login.V)(vars).CurrentUser()
	r := (httper.V)(vars).GetRequest()

	t := trigger.New()

	sendMessage := func(message string) {
		var jsonMessage interface{}
		err := json.Unmarshal([]byte(message), &jsonMessage)
		if err != nil {
			(logger.V)(vars).Get().Println(message, "is not valid JSON")
			return
		}
		(logger.V)(vars).Get().Println(message, "gets translated into", jsonMessage)
		(websocketer.V)(vars).SendJSONMessage(jsonMessage)
	}

	if loggedIn {
		loginChannel(currentUser.ClashTag).Publish("Doesn't matter what goes here")
		closer := pubsuber.User(currentUser.ClashTag).ReceiveMessages(sendMessage)
		t.OnClose(func() {
			closer.Close()
		})
	}
	closer := pubsuber.Url(r.URL.String()).ReceiveMessages(sendMessage)
	t.OnClose(func() {
		closer.Close()
	})

	vars[closerIndex] = t
}

var wsCloser rack.Func = func(vars map[string]interface{}, next func()) {
	closer := vars[closerIndex].(io.Closer)
	closer.Close()

	currentUser, loggedIn := (login.V)(vars).CurrentUser()
	if loggedIn {
		if !pubsuber.User(currentUser.ClashTag).IsActive() {
			go startLogoutProcess(currentUser.ClashTag)
		}
	}
}

var wsMessager rack.Func = func(vars map[string]interface{}, next func()) {

	message, ok := (websocketer.V)(vars).GetMessage().(*WebsocketMessage)
	if !ok {
		(logger.V)(vars).Get().Println("Unknown Message")
		return
	}

	actionType := message.Type
	actionData := message.Data

	action, ok := messageTypes[actionType]
	if !ok {
		(logger.V)(vars).Get().Println("No action for messages of type", actionType)
		return
	}

	currentUser, _ := (login.V)(vars).CurrentUser()

	result := action(currentUser, actionData)

	(websocketer.V)(vars).SetResponse(result)
}

var messageTypes = make(map[string]func(*user.User, interface{}) interface{})

func MessageAction(messagetype string, action func(*user.User, interface{}) interface{}) {
	messageTypes[messagetype] = action
}
