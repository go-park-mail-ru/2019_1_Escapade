package models

import "time"

// Game show all info about game room.
// There is no any personal info about gamer
type Game struct {
	RoomID        string    `json:"roomID"`
	Name          string    `json:"name"`
	Status        int       `json:"status"`
	Players       int       `json:"players"`
	TimeToPrepare int       `json:"prepare"`
	TimeToPlay    int       `json:"play"`
	Date          time.Time `json:"date"`
}

// Gamer show all personal info(gamers results) about game
type Gamer struct {
	ID         int     `json:"-"`
	Score      float64 `json:"score"`
	Time       int     `json:"time"`
	LeftClick  int     `json:"leftClick"`
	RightClick int     `json:"rightClick"`
	Explosion  bool    `json:"online"`
	Won        bool    `json:"won"`
}

// Action is the database model of game.Action
type Action struct {
	PlayerID int       `json:"playerID"`
	ActionID int       `json:"actionID"`
	Date     time.Time `json:"-"`
}

// Cell is the database model of game.Cell
type Cell struct {
	PlayerID int       `json:"playerID"`
	X        int       `json:"x"`
	Y        int       `json:"y"`
	Value    int       `json:"value"`
	Date     time.Time `json:"-"`
}

// Field is the database model of game.Field
type Field struct {
	Width     int `json:"width"`
	Height    int `json:"height"`
	CellsLeft int `json:"cellsLeft"`
	Difficult int `json:"difficult"`
	Mines     int `json:"mines"`
}

// GameInformation show everything about game and his gamer
type GameInformation struct {
	Game     Game       `json:"game"`
	Field    Field      `json:"field"`
	Actions  []Action   `json:"actions"`
	Cells    []Cell     `json:"cells"`
	Gamers   []Gamer    `json:"gamer"`
	Messages []*Message `json:"messages"`
}
