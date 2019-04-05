package game

import (
	"escapade/internal/models"
	//re "escapade/internal/return_errors"
	//"math/rand"
)

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
	Cell   *models.Cell `json:"cell"`
	Action *int         `json:"action"`
}

type RoomGet struct {
	Players   bool `json:"players"`
	Observers bool `json:"observers"`
	Field     bool `json:"field"`
	History   bool `json:"history"`
}
