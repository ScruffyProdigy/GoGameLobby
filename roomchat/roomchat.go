package roomchat

import (
	"../models/user"
	"../pubsuber"
	"../websocketcontrol"
	"time"
)

func init() {
	websocketcontrol.MessageAction("roomchat", func(u *user.User, i_message interface{}) interface{} {
		m_message := i_message.(map[string]interface{})
		url := m_message["url"].(string)
		text := m_message["text"].(string)
		type roomchatinfo struct {
			Name string `json:'name'`
			Text string `json:'text'`
			Time string `json:'time'`
		}
		print("Sending Message\n")
		pubsuber.Url(url).SendMessage("roomchat", roomchatinfo{Name: u.ClashTag, Text: text, Time: time.Now().Format("2006-01-02T15:04:05-07:00")})
		return nil
	})
}
