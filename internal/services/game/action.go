package game

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
)

// PlayerAction combine player and his action
type PlayerAction struct {
	Player int `json:"player"`
	Action int `json:"action"`
}

// NewPlayerAction return new instance of PlayerAction
func NewPlayerAction(player int, action int) *PlayerAction {
	pa := &PlayerAction{
		Player: player,
		Action: action,
	}
	return pa
}

// Free free memory
func (pa *PlayerAction) Free() {
	if pa == nil {
		return
	}
	pa = nil
}

func (room *Room) addAction(conn *Connection, action int) {
	pa := NewPlayerAction(conn.ID(), action)
	room.History = append(room.History, pa)
}
