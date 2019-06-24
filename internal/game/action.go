package game

import "time"

// Player actions
const (
	ActionError             = -1
	ActionNo                = 0
	ActionConnectAsPlayer   = 1
	ActionConnectAsObserver = 2
	ActionReconnect         = 3
	ActionDisconnect        = 4
	//ActionDisconnectObserver = 5
	ActionFlagСonflict = 6
	ActionExplode      = 7
	ActionWin          = 8
	ActionLose         = 9
	ActionFlagLost     = 10
	ActionGetPoints    = 11
	ActionFlagSet      = 12
	ActionGiveUp       = 13
	ActionBackToLobby  = 14
	ActionTimeOver     = 15
	ActionRestart      = 16
	ActionTimeout      = 17
	ActionConnect      = 18
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
