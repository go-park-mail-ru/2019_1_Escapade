package game

import "time"

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
	ActionBackToLobby
	ActionTimeOver
)

// PlayerAction combine player and his action
type PlayerAction struct {
	Player int       `json:"player"`
	Action int       `json:"action"`
	Time   time.Time `json:"-"`
}

// NewPlayerAction return new instance of PlayerAction
func NewPlayerAction(player int, action int) *PlayerAction {
	pa := &PlayerAction{
		Player: player,
		Action: action,
		Time:   time.Now(),
	}
	return pa
}
