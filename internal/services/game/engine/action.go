package engine

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
//easyjson:json
type PlayerAction struct {
	Player int32     `json:"player"`
	Action int32     `json:"action"`
	Time   time.Time `json:"-"`
}

// NewPlayerAction return new instance of PlayerAction
func NewPlayerAction(player int32, action int32) *PlayerAction {
	pa := &PlayerAction{
		Player: player,
		Action: action,
		Time:   time.Now(),
	}
	return pa
}
