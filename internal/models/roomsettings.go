package models

type RoomSettings struct {
	id      int
	Width   int
	Height  int
	Players int
	Percent float32
}

func NewSmallRoom() *RoomSettings {
	rs := &RoomSettings{
		Width:   7,
		Height:  7,
		Players: 2,
		Percent: 0.15,
	}
	return rs
}

func NewUsualRoom() *RoomSettings {
	rs := &RoomSettings{
		Width:   12,
		Height:  12,
		Players: 2,
		Percent: 0.17,
	}
	return rs
}

func NewBigRoom() *RoomSettings {
	rs := &RoomSettings{
		Width:   20,
		Height:  20,
		Players: 2,
		Percent: 0.25,
	}
	return rs
}
