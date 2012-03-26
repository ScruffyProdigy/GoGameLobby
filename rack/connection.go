package rack

import "net/http"

type Connection interface {
	Go(func(http.ResponseWriter, *http.Request))
}

type httpConnection struct {
	address string
}

func (this *httpConnection) Go(f func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc("/", f)
	http.ListenAndServe(this.address, nil)
}

func HttpConnection(addr string) Connection {
	conn := new(httpConnection)
	conn.address = addr
	return conn
}

type httpsConnection struct {
	address  string
	certFile string
	keyFile  string
}

func (this *httpsConnection) Go(f func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc("/", f)
	http.ListenAndServeTLS(this.address, this.certFile, this.keyFile, nil)
}

func HttpsConnection(addr, certFile, keyFile string) Connection {
	conn := new(httpsConnection)
	conn.address = addr
	conn.certFile = certFile
	conn.keyFile = keyFile

	return conn
}
