package constants

const (
	Debug = iota
	Release
)

const (
	Mode = Debug
)

var Base string
var Port string
var Site string

func init() {
	if Mode == Debug {
		Base = "http://localhost"
		Port = ":3000"
		Site = Base + Port
	} else if Mode == Release {
		Base = "http://clashcentral.com"
		Port = ":80"
		Site = Base
	}
}
