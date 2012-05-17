package messenger

import (
	"bytes"
	"encoding/json"
	"io"
)

func JSONmessage(message interface{}, url string, response interface{}) error {
	m, err := jsoncodec.CreateMessage(message)
	if err != nil {
		return err
	}

	return m.GetResponse(url, response)
}

func JSONmessageNew(message interface{}, url string) (result interface{}, err error) {
	err = JSONmessage(message, url, &result)
	return
}

var jsoncodec = Codec{
	Mime: "application/json",
	Encode: func(original interface{}) (io.Reader, error) {
		b, err := json.Marshal(original)
		return bytes.NewBuffer(b), err
	},
	Decode: func(r io.Reader, result interface{}) error {
		d := json.NewDecoder(r)
		return d.Decode(result)
	},
}

func init() {
	RegisterCodec(&jsoncodec)
}
