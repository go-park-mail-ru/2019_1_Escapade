package game

/* Examples of json

room search
{"send":{"RoomSettings":{"name":"","width":12,"height":12,"players":2,"observers":10,"mines":5}},"get":null}

send cell
{"send":{"cell":{"x":2,"y":1,"value":0,"PlayerID":0}, "action":null},"get":null}

send action(all actions are in action.go). Server iswaiting only one of these:
ActionStop 5
ActionContinue 6
ActionGiveUp 13
ActionBackToLobby 14

give up
{"send":{"cell":null, "action":13,"get":null}

back to lobby
{"send":{"cell":null, "action":14,"get":null}

get lobby all info
{"send":null,"get":{"allRooms":true,"freeRooms":true,"waiting":true,"playing":true}}


*/

type RoomRequest struct {
	Send *RoomSend `json:"send"`
	Get  *RoomGet  `json:"get"`
}

func (rr *RoomRequest) IsGet() bool {
	return rr.Get != nil
}

func (rr *RoomRequest) IsSend() bool {
	return rr.Send != nil
}

type RoomSend struct {
	Cell   *Cell `json:"cell,omitempty"`
	Action *int  `json:"action,omitempty"`
}

type RoomGet struct {
	Players   bool `json:"players"`
	Observers bool `json:"observers"`
	Field     bool `json:"field"`
	History   bool `json:"history"`
}
