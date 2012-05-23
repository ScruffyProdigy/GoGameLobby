package websocketcontrol

import (
	"../../Middleware/websocketer"
	"../login"
	"../models/user"
	"../pubsuber"
	"../redis"
	"../trigger"
	"encoding/json"
	"fmt"
	"github.com/HairyMezican/Middleware/logger"
	"github.com/HairyMezican/TheRack/httper"
	"github.com/HairyMezican/TheRack/rack"
	"io"
	"time"
)

const (
	closerIndex = "websocket:closer"
)

func cancelLogoutIndex(user string) string {
	return "Cancel Logout " + user
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
		fmt.Println("Making New Storage")
		var m WebsocketMessage
		return &m
	})
	fmt.Println("I Did OnStorage")

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
		redis.Client.Publish(cancelLogoutIndex(currentUser.ClashTag), "Doesn't matter what goes here")
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
