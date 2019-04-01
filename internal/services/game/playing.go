package game

import "escapade/internal/models"

// Playing is a player, that playing game
type Playing struct {
	Player   *models.Player
	Flag     *models.Flag
	InGame   bool
	Finished bool
	Stoped   bool // player waits
}

// NewPlaying create instance of Playing before game starts
func NewPlaying(player *models.Player, flag *models.Flag) *Playing {
	playing := &Playing{
		Player:   player,
		Flag:     flag,
		InGame:   true,
		Finished: false,
		Stoped:   false,
	}
	return playing
}
