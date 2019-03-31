package game

import "escapade/internal/models"

// Playing is a player, that playing game
type Playing struct {
	Player      *Player
	Flag        *models.Flag
	InGame      bool
	WasCaptured bool
	Stoped      bool // player waits
	TakenFlags  int
	points      int
}

// NewPlaying create instance of Playing before game starts
func NewPlaying(player *Player, flag *models.Flag) *Playing {
	playing := &Playing{
		Player:      player,
		Flag:        flag,
		InGame:      true,
		WasCaptured: false,
		TakenFlags:  0,
		points:      0,
	}
	return playing
}
