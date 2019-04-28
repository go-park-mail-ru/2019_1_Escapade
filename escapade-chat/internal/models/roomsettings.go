package models

// RoomSettings a set of parameters, that set the size,  of room
// and field, the duration of game etc
type RoomSettings struct {
	Name          string `json:"name"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	Players       int    `json:"players"`
	Observers     int    `json:"observers"`
	TimeToPrepare int    `json:"prepare"`
	TimeToPlay    int    `json:"play"`
	Mines         int    `json:"mines"`
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
