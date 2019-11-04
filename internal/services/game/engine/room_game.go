package engine

// IsWinner is player wuth id playerID is winner

///////// sent to onlinePlayers
/*
func (room *Room) isWinner(playerIndex int, isMax *bool) func(int, Player) {
	var (
		max              = 0.
		thisPlayerPoints float64
		ignore           bool
	)

	return func(index int, player Player) {
		if ignore || player.Died || player.Points < max {
			return
		}
		max = player.Points
		if index == playerIndex {
			thisPlayerPoints = max
			*isMax = true
		} else if *isMax && max > thisPlayerPoints {
			*isMax = false
			ignore = true
		}
	}
}*/
