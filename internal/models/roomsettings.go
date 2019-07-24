package models

// RoomSettings a set of parameters, that set the size,  of room
// and field, the duration of game etc
//easyjson:json
type RoomSettings struct {
	ID            string `json:"id"`
	Name          string `json:"name" minLength:"1" maxLength:"30"`
	Width         int    `json:"width" minimum:"5" maximum:"100"`
	Height        int    `json:"height" minimum:"5" maximum:"100"`
	Players       int    `json:"players" minimum:"2" maximum:"100"`
	Observers     int    `json:"observers"  maximum:"100"`
	TimeToPrepare int    `json:"prepare" maximum:"60"`
	TimeToPlay    int    `json:"play" maximum:"7200"`
	Mines         int    `json:"mines"`
	NoAnonymous   bool   `json:"noAnonymous"`
	Deathmatch    bool   `json:"deathmatch"`
}

// FieldCheck check room's field is valid
func (rs *RoomSettings) FieldCheck() bool {
	return rs.Players > 1 && rs.Width*rs.Height-rs.Players-rs.Mines > 0
}

// AnonymousCheck check that anonymous cant join to non anonymous game
func (rs *RoomSettings) AnonymousCheck(isAnonymous bool) bool {
	return !rs.NoAnonymous || !isAnonymous
}

// Similar compare two Rooms settings
func (rs *RoomSettings) Similar(another *RoomSettings) bool {
	return rs.Width == another.Width && rs.Height == another.Height &&
		rs.TimeToPlay == another.TimeToPlay && rs.Mines == another.Mines &&
		rs.TimeToPrepare == another.TimeToPrepare && rs.Players == another.Players &&
		rs.Observers == another.Observers && rs.NoAnonymous == another.NoAnonymous &&
		rs.Deathmatch == another.Deathmatch
}

// NewSmallRoom create small room
func NewSmallRoom() *RoomSettings {
	rs := &RoomSettings{
		Width:         7,
		Height:        7,
		Players:       2,
		Observers:     10,
		TimeToPrepare: 5,
		TimeToPlay:    60,
		Mines:         2,
	}
	return rs
}

// NewUsualRoom create middle room
func NewUsualRoom() *RoomSettings {
	rs := &RoomSettings{
		Width:         12,
		Height:        12,
		Players:       2,
		Observers:     10,
		TimeToPrepare: 10,
		TimeToPlay:    120,
		Mines:         5,
	}
	return rs
}

// NewBigRoom create big room
func NewBigRoom() *RoomSettings {
	rs := &RoomSettings{
		Width:         20,
		Height:        20,
		Players:       2,
		Observers:     10,
		TimeToPrepare: 15,
		TimeToPlay:    300,
		Mines:         10,
	}
	return rs
}
