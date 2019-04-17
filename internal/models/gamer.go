package models

// Gamer show all personal info(gamers results) about game
type Gamer struct {
	ID         int  `json:"-"`
	Score      int  `json:"score"`
	Time       int  `json:"time"`
	LeftClick  int  `json:"leftClick"`
	RightClick int  `json:"rightClick"`
	Explosion  bool `json:"online"`
	Won        bool `json:"won"`
}
