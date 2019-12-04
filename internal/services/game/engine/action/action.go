package action

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
)

// Player actions
const (
	Error             = -1
	No                = 0
	ConnectAsPlayer   = 1
	ConnectAsObserver = 2
	Reconnect         = 3
	Disconnect        = 4
	//ActionDisconnectObserver = 5
	FlagСonflict = 6
	Explode      = 7
	Win          = 8
	Lose         = 9
	FlagLost     = 10
	GetPoints    = 11
	FlagSet      = 12
	GiveUp       = 13
	BackToLobby  = 14
	TimeOver     = 15
	Restart      = 16
	Timeout      = 17
	Connect      = 18
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

// ToModel cast PlayerAction to models.Action
func (action *PlayerAction) ToModel() models.Action {
	return models.Action{
		PlayerID: action.Player,
		ActionID: action.Action,
		Date:     action.Time,
	}
}

// FromModel cast models.Action to PlayerAction
func (action *PlayerAction) FromModel(model models.Action) *PlayerAction {
	action.Player = model.PlayerID
	action.Action = model.ActionID
	action.Time = model.Date
	return action
}
