package redis

import (
	"errors"
	"github.com/simonz05/godis/redis"
)

type Subscription chan<- bool

func (this *Subscription) Close() error {
	if *this == nil {
		return errors.New("Already closed this subscription")
	}
	*this <- true
	close(*this)
	*this = nil
	return nil
}

func waitForMessages(subscriber *redis.Sub, action func(string), closer func()) *Subscription {
	trigger := make(chan bool)
	subscription := Subscription(trigger)
	go func() {
		defer closer()

		for {
			select {
			case m := <-subscriber.Messages:
				action(m.Elem.String())
			case <-trigger:
				return
			}
		}
	}()
	return &subscription
}

func (this Channel) Subscribe(action func(string)) *Subscription {
	subscriber, err := this.client.newClient().Subscribe(this.key)
	checkForError(err)
	return waitForMessages(subscriber, action, func() {
		subscriber.Unsubscribe(this.key)
		subscriber.Close()
	})
}

func (this Channel) PatternSubscribe(action func(string)) *Subscription {
	subscriber, err := this.client.newClient().Psubscribe(this.key)
	checkForError(err)
	return waitForMessages(subscriber, action, func() {
		subscriber.Punsubscribe(this.key)
		subscriber.Close()
	})
}

func translator(elements <-chan *redis.Message, strings chan<- string) {
	for message := range elements {
		strings <- message.Elem.String()
	}
}

func (this Channel) BlockingSubscription(subscription func(<-chan string)) {
	subscriber, err := this.client.newClient().Subscribe(this.key)
	checkForError(err)

	stringChannel := make(chan string)
	go translator(subscriber.Messages, stringChannel)
	subscription(stringChannel)

	subscriber.Unsubscribe(this.key)
	subscriber.Close()
}

func (this Channel) BlockingPatternSubscription(subscription func(<-chan string)) {
	subscriber, err := this.client.newClient().Psubscribe(this.key)
	checkForError(err)

	stringChannel := make(chan string)
	go translator(subscriber.Messages, stringChannel)
	subscription(stringChannel)

	subscriber.Punsubscribe(this.key)
	subscriber.Close()
}

func (this Channel) Publish(message string) int64 {
	receivers, err := this.client.client.Publish(this.key, message)
	checkForError(err)
	return receivers
}

type Channel struct {
	client Redis
	key    string
}

func newChannel(client Redis, key string) Channel {
	return Channel{
		client: client,
		key:    key,
	}
}
