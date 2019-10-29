package game

// Player has name, ID, points and flag Finish
type Player struct {
	ID       int32
	Points   float64
	Finished bool
	Died     bool
}

// NewPlayer create new instance of Player
func NewPlayer(id int32) *Player {
	player := &Player{
		ID:     id,
		Points: 0,
	}
	return player
}

// SetAsPlayer set Finished = false
func (player *Player) SetAsPlayer() {
	player.Points = 0
	player.Finished = false
}

// SetAsObserver set Finished = true
func (player *Player) SetAsObserver() {
	player.Points = 0
	player.Finished = true
}

/*
During the game the room doesnt know status of that player
(is he gamer or observer). Every player can send cell and
action. But the room process these requests only from players,
which 'finish' flag is false(except action - Back to menu)
*/

// IsAlive check is player alive
func (player *Player) IsAlive() bool {
	return player.Finished == false
}
