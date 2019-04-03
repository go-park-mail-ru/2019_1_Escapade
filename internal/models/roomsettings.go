package models

type RoomSettings struct {
	Name    string `json:"name"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Players int    `json:"players"`
	Mines   int    `json:"mines"`
}

func NewSmallRoom() *RoomSettings {
	rs := &RoomSettings{
		Width:   7,
		Height:  7,
		Players: 2,
		Mines:   2,
	}
	return rs
}

func NewUsualRoom() *RoomSettings {
	rs := &RoomSettings{
		Width:   12,
		Height:  12,
		Players: 2,
		Mines:   5,
	}
	return rs
}

func NewBigRoom() *RoomSettings {
	rs := &RoomSettings{
		Width:   20,
		Height:  20,
		Players: 2,
		Mines:   10,
	}
	return rs
}
