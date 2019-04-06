package game

//re "escapade/internal/return_errors"
//"math/rand"

// Player actions
const (
	ActionError = iota - 1
	ActionNo
	ActionConnectAsPlayer
	ActionConnectAsObserver
	ActionReconnect
	ActionDisconnect
	ActionStop
	ActionContinue
	ActionExplode
	ActionWin
	ActionLose
	ActionFlagLost
	ActionGetPoints
	ActionFlagSet
	ActionGiveUp
)

// PlayerAction combine player and his action
type PlayerAction struct {
	Player *Player `json:"player"`
	Action int     `json:"action"`
}

// NewPlayerAction return new instance of PlayerAction
func NewPlayerAction(player *Player, action int) *PlayerAction {
	pa := &PlayerAction{
		Player: player,
		Action: action,
	}
	return pa
}

func (room *Room) addAction(conn *Connection, action int) {
	pa := NewPlayerAction(conn.Player, action)
	room.History = append(room.History, pa)
}
