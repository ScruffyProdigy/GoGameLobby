package trigger

import "errors"

type Trigger struct {
	pulled  bool
	onclose []func()
}

func New() *Trigger {
	t := new(Trigger)
	t.onclose = make([]func(), 0)
	return t
}

func OnClose(action func()) *Trigger {
	t := New()
	t.OnClose(action)
	return t
}

func (this *Trigger) Close() error {
	if this.pulled {
		return errors.New("Trigger Already Pulled")
	}

	this.pulled = true
	for _, action := range this.onclose {
		action()
	}
	return nil
}

func (this *Trigger) OnClose(action func()) {
	this.onclose = append(this.onclose, action)
}

func (this *Trigger) Channel() <-chan bool {
	closer := make(chan bool)
	this.OnClose(func() {
		go func() { closer <- true }()
	})
	return closer
}
