package models

type PlayerStatistics struct {
	Name        string `json:"name"`
	GamesTotal  int    `json:"photo"`
	SingleTotal int    `json:"bestScore"`
	OnlineTotal int    `json:"bestTime"`
	SingleWin   int    `json:"singleWin"`
	OnlineWin   int    `json:"onlineWin"`
	MinsFound   int    `json:"minsFound"`
	FirstSeen   string `json:"firstSeen"`
	LastSeen    string `json:"lastSeen"`
}
