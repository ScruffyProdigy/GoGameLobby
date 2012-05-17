package gamedata

//
//  Mode List Querying - Which Modes can this user current select? (and which groups within those modes can he select?)
//
type PreModeInfo struct {
	User string `json:"user"` //The name of the user within the lobby
}

type ModeInfo struct {
	Name   string                       `json:"name"`   //a human readable description of the game mode
	Groups *map[string]GroupDescription `json:"groups"` //a map listing the groups within the mode that the player is allowed to join
}

type GroupDescription struct {
	Name string `json:"name"` //a human readable name for the group
}

//
//  Join Querying - What player options can this user currently select when joining the game?
//
type PreJoinInfo struct {
	User  string `json:"user"`  //the name of the user from the lobby
	Mode  string `json:"mode"`  //the machine readable descriptor you sent in the clash info ("Identifier" in ModeInfo)
	Group string `json:"group"` //the group that the player wants to join (from "Players" in ModeInfo)
}

type JoinInfo struct {
	Name string `json:"name"` // a human readable description
}

//
//Player Counting - How many players does the lobby need in each group before it can send you the players for a game?
//
type PrePlayerCountInfo struct {
	Mode string `json:"mode"` // the machine-readable name of the mode that they want to get the player counts for
}

type ModePlayerCounts struct {
	Players map[string]int `json:"players"` //key: the machine-readable name of the group; value: the number of people that fit in the group
}

//
//Start a Clash! - Here are the players that want to clash!
//
type PreStartInfo struct {
	Mode      string                   `json:"mode"`   //"Identifier" from ModeInfo
	Groups    map[string]GroupJoinInfo `json:"groups"` //key: each of the groups that you specified in both ModePlayersCounts and ModeInfo
	ResultUrl string                   `json:"result"` //the url to POST/PUT the results to when the clash is finished
}

type GroupJoinInfo struct {
	Players map[string]PlayerJoinInfo `json:"players"` //key: the username of each player in the group
}

type PlayerJoinInfo struct {
	Join string `json:"Join"` //the index into the JoinInfo map that the player chose
}

type StartInfo struct {
	Url     string                     `json:"url"`     //the url to send all of the players to
	Players map[string]PlayerStartInfo `json:"players"` //key: the clashtag of each user in the game
}

type PlayerStartInfo struct {
	UrlValues map[string]string `json:"urlvals"` //the URL values to send to each of the players (usually used so you can identify each player when they get to the game)
}
