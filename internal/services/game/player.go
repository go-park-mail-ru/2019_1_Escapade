package game

type Player struct {
	ID     int
	Name   string
	Points int
}

func NewPlayer(name string, id int) *Player {
	player := &Player{
		ID:     id,
		Name:   name,
		Points: 0,
	}
	return player
}

// Reset - call it in every game beginning
func (player *Player) Reset() {
	player.Points = 0
}
