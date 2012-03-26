package routes

import "net/http"
import "strings"
import "fmt"

type Urler interface {
	Url() []string
}

func RedirectWithCode(this http.ResponseWriter, code int, redirecttome Urler) {
	url := "/" + strings.Join(redirecttome.Url(), "/")
	this.Header().Add("Location", url)
	this.WriteHeader(code)
	fmt.Fprint(this, "You are being redirected to ", url)
}

func RedirectTo(this http.ResponseWriter, redirectome Urler) {
	RedirectWithCode(this, http.StatusFound, redirectome)
}
