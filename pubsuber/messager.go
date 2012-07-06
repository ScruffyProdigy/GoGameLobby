package pubsuber

import (
	"github.com/HairyMezican/SimpleRedis/redis"
	"../global"
	"../trigger"
	"encoding/json"
	"io"
)

func userMessageChannel(user string) redis.Channel {
	return global.Redis.Channel("User Channel:" + user)
}

func userStoredMessages(user string) redis.List {
	return global.Redis.List("Users Stored Messages:" + user)
}

func urlMessageChannel(url string) redis.Channel {
	return global.Redis.Channel("Url Channel:" + url)
}

func userInstanceCount(user string) redis.Integer {
	return global.Redis.Integer("User Instances " + user)
}

type Target interface {
	SendMessage(message interface{})
	ReceiveMessages(action func(message string)) io.Closer
	IsActive() bool
}

type UserTarget interface {
}

func User(name string) Target {
	return userTarget{name}
}

type userTarget struct {
	user string
}

func makeString(message interface{}) string {

	stringMessage, ok := message.(string)
	if ok {
		return stringMessage
	}

	byteMessage, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	return string(byteMessage)
}

func (this userTarget) SendMessage(message interface{}) {
	stringMessage := makeString(message)
	if this.IsActive() {
		userMessageChannel(this.user).Publish(stringMessage)
	} else {
		userStoredMessages(this.user).LeftPush(stringMessage)
	}
}

func (this userTarget) ReceiveMessages(action func(message string)) io.Closer {
	_,subscription := userMessageChannel(this.user).Subscribe(action)
	userInstanceCount(this.user).Increment()

	for {
		message, ok := <-userStoredMessages(this.user).RightPop()
		if !ok {
			break
		}
		action(message)
	}

	return trigger.OnClose(func() {
		subscription.Close()
		userInstanceCount(this.user).Decrement()
	})
}

func (this userTarget) IsActive() bool {
	return <-userInstanceCount(this.user).Get() > 0
}

func Url(url string) Target {
	return urlTarget{url}
}

type urlTarget struct {
	url string
}

func (this urlTarget) SendMessage(message interface{}) {
	stringMessage := makeString(message)
	urlMessageChannel(this.url).Publish(stringMessage)
}

func (this urlTarget) ReceiveMessages(action func(message string)) io.Closer {
	_,result := urlMessageChannel(this.url).Subscribe(action)
	return result
}

func (this urlTarget) IsActive() bool {
	return true
}
