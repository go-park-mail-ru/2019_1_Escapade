package models

import "time"

// Game show all info about game room.
// There is no any personal info about gamer
//easyjson:json
type Game struct {
	ID              int32         `json:"-"`
	Settings        *RoomSettings `json:"settings"`
	Status          int32         `json:"status"`
	RecruitmentTime time.Duration `json:"recruitmentTime"`
	PlayingTime     time.Duration `json:"playingTime"`
	Date            time.Time     `json:"date"`
	ChatID          int32         `json:"-"`
}

// Gamer show all personal info(gamers results) about game
//easyjson:json
type Gamer struct {
	ID         int32   `json:"-"`
	Score      float64 `json:"score"`
	Time       int32   `json:"time"`
	LeftClick  int32   `json:"leftClick"`
	RightClick int32   `json:"rightClick"`
	Explosion  bool    `json:"online"`
	Won        bool    `json:"won"`
}

// Action is the database model of game.Action
//easyjson:json
type Action struct {
	PlayerID int32     `json:"playerID"`
	ActionID int32     `json:"actionID"`
	Date     time.Time `json:"-"`
}

// Cell is the database model of game.Cell
//easyjson:json
type Cell struct {
	PlayerID int32     `json:"playerID"`
	X        int32     `json:"x"`
	Y        int32     `json:"y"`
	Value    int32     `json:"value"`
	Date     time.Time `json:"-"`
}

// Field is the database model of game.Field
//easyjson:json
type Field struct {
	Width     int32 `json:"width"`
	Height    int32 `json:"height"`
	CellsLeft int32 `json:"cellsLeft"`
	Difficult int   `json:"difficult"`
	Mines     int32 `json:"mines"`
}

// GameInformation show everything about game and his gamer
//easyjson:json
type GameInformation struct {
	Game     Game       `json:"game"`
	Field    Field      `json:"field"`
	Actions  []Action   `json:"actions"`
	Cells    []Cell     `json:"cells"`
	Gamers   []Gamer    `json:"gamer"`
	Messages []*Message `json:"messages"`
}
