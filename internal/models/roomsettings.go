package models

// RoomSettings a set of parameters, that set the size,  of room
// and field, the duration of game etc
type RoomSettings struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	Players       int    `json:"players"`
	Observers     int    `json:"observers"`
	TimeToPrepare int    `json:"prepare"`
	TimeToPlay    int    `json:"play"`
	Mines         int    `json:"mines"`
	NoAnonymous   bool   `json:"noAnonymous"`
	Deathmatch    bool   `json:"deathmatch"`
}

func (rs *RoomSettings) FieldCheck() bool {
	return rs.Players > 1 && rs.Width*rs.Height-rs.Players-rs.Mines > 0
}

func (rs *RoomSettings) AnonymousCheck(isAnonymous bool) bool {
	return !rs.NoAnonymous || !isAnonymous
}

func (origin *RoomSettings) Similar(another *RoomSettings) bool {
	return origin.Width == another.Width && origin.Height == another.Height &&
		origin.TimeToPlay == another.TimeToPlay && origin.Mines == another.Mines &&
		origin.TimeToPrepare == another.TimeToPrepare && origin.Players == another.Players &&
		origin.Observers == another.Observers && origin.NoAnonymous == another.NoAnonymous &&
		origin.Deathmatch == another.Deathmatch
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
