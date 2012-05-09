package gamedata

type PreModeInfo struct {
	User string `json:"user"` //The name of the user within the lobby
}

type Mode struct {
	Identifier string         `json:"identifier"` //a machine readable description of the game mode (this is what you'll receive back)
	Players    map[string]int `json:"players"`    //a list of all of the player groups that need to be filled, and how many players it will take to fill them
}

type ModeInfo struct {
	Descriptor string `json:"descriptor"` //a human readable description of the game mode
	Mode       *Mode  `json:"mode"`       //the information about the mode
}

type PreJoinInfo struct {
	User  string `json:"user"`  //the name of the user from the lobby
	Mode  string `json:"mode"`  //the machine readable descriptor you sent in the clash info ("Identifier" in ModeInfo)
	Group string `json:"group"` //the group that the player wants to join (from "Players" in ModeInfo)
}

type Join struct {
	PublicOptions map[string]string //Human readable description
	Identifier    string            // Machine readble description
}

type JoinInfo struct {
	Description string
	Join        *Join
}

type PreStartInfo struct {
	Mode      string                       `json:"mode"`    //"Identifier" from ModeInfo
	Players   map[string]map[string]string `json:"players"` //similar to "Players" from ModeInfo, except replacing the integer with a map (with a len equal to the integer) linking the user's lobby name to their JoinInfo Identifier
	ResultUrl string                       `json:"result"`  //the url to POST/PUT the results to when the clash is finished
}

type StartInfo struct {
	Url  string
	Info map[string]map[string]string //a map which links a user's lobby name to the options that should be sent in when sending that user to the clash
}
