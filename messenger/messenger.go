package messenger

import (
	"io"
	"net/http"
)

type Codec struct {
	Mime   string
	Encode func(interface{}) (io.Reader, error)
	Decode func(io.Reader, interface{}) error
}

type codecMap map[string]Codec

type message struct {
	message io.Reader
	mime    string
}

var mimeTypes codecMap = make(codecMap)

func RegisterCodec(c Codec) {
	mimeTypes[c.Mime] = c
}

type messageHandler func(m message) error

func (this Codec) CreateMessage(content interface{}) (*message, error) {
	r, err := this.Encode(content)
	if err != nil {
		return nil, err
	}
	result := new(message)
	result.message = r
	result.mime = this.Mime
	return result, nil
}

func (this codecMap) DecodeMessage(m message, result interface{}) error {
	codec := this[m.mime]
	return codec.Decode(m.message, result)
}

func (this message) Read(p []byte) (int, error) {
	return this.message.Read(p)
}

func (this message) SendTo(url string, mh messageHandler) error {
	r, err := http.Post(url, this.mime, this)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return mh(message{
		mime:    r.Header.Get("content-type"),
		message: r.Body,
	})
}

func (this message) GetResponse(url string, response interface{}) error {
	return this.SendTo(url, func(m message) error {
		return mimeTypes.DecodeMessage(m, response)
	})
}
