package rack

import (
	"net/http"
)

type fake struct {
	status  int
	message []byte
	header  http.Header
}

type FakeResponseWriter interface {
	http.ResponseWriter
	Results() (int, http.Header, []byte)
}

func BlankResponse() FakeResponseWriter {
	this := new(fake)
	this.status = http.StatusOK
	this.header = make(http.Header)
	return this
}

func CreateResponse(status int, header http.Header, message []byte) FakeResponseWriter {
	this := new(fake)
	this.status = status
	this.header = header
	this.message = message
	return this
}

func (this *fake) Header() http.Header {
	return this.header
}

func (this *fake) Write(message []byte) (bytes int, err error) {
	this.message = append(this.message, message...)
	bytes += len(message)
	return
}

func (this *fake) WriteHeader(status int) {
	this.status = status
}

func (this *fake) Results() (int, http.Header, []byte) {
	return this.status, this.header, this.message
}
