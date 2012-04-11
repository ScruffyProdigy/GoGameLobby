/*
	session wraps gorilla sessions within a Rack Middleware framework
*/
package session

import (
	"code.google.com/p/gorilla/sessions"
	"net/http"
)

var store = sessions.NewCookieStore([]byte("Go Game Lobby!"))

/*
	The Session is the interface exposed to the rest of the program
*/
type Session interface {
	Set(k, v interface{})          // Set will set a session variable
	Get(k interface{}) interface{} //Get will obtain the result of a session variable
	Clear(k interface{})           //Clear will get rid of a session variable
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

func (this *session) Clear(k interface{}) {
	delete(this.sess.Values, k)
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
