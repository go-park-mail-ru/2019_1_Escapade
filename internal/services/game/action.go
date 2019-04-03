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
	ActionGetPoints
	ActionFlagSet
	ActionGiveUp
)

type PlayerAction struct {
	Player *Player `json:"player"`
	Action int     `json:"action"`
}

func NewPlayerAction(player *Player, action int) *PlayerAction {
	pa := &PlayerAction{
		Player: player,
		Action: action,
	}
	return pa
}
