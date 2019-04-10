package models

import "time"

// Game show all info about game room.
// There is no any personal info about gamer
type Game struct {
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Players   int       `json:"players"`
	Mines     int       `json:"mines"`
	Date      time.Time `json:"date"`
	Online    bool      `json:"online"`
	Difficult int       `json:"difficult"`
}

// GameInformation show everything about game and his gamer
type GameInformation struct {
	Game  *Game  `json:"game"`
	Gamer *Gamer `json:"gamer"`
}
