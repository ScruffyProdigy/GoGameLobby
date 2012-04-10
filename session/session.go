package session

import (
	"code.google.com/p/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("Go Game Lobby!"))

type Session interface {
	Set(k, v interface{})
	Get(k interface{}) interface{}
}

type session struct {
	sess *sessions.Session
	r    *http.Request
}

func (this *session) Set(k, v interface{}) {
	this.sess.Values[k] = v
}

func (this *session) Get(k interface{}) interface{} {
	return this.sess.Values[k]
}

func (this *session) save(w http.ResponseWriter) {
	this.sess.Save(this.r, w)
}

func get(r *http.Request) *session {
	sess := new(session)
	sess.sess, _ = store.Get(r, "session")
	sess.r = r
	return sess
}
