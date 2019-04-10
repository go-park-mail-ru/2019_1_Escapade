package models

import "time"

// Gamer show all personal info(gamers results) about game
type Gamer struct {
	Score      int       `json:"score"`
	Time       time.Time `json:"time"`
	MinesOpen  int       `json:"minesOpen"`
	LeftClick  int       `json:"leftClick"`
	RightClick int       `json:"rightClick"`
	Explosion  bool      `json:"online"`
	Won        bool      `json:"won"`
}
