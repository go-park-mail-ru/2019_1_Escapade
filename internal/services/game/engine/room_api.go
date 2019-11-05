package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// APIStrategyI handle client requests to room
// Strategy Pattern
type APIStrategyI interface {
	Handle(conn *Connection, rr *RoomRequest)
}

// RoomAPI implement APIStrategyI
type RoomAPI struct {
	s  SyncI
	m  MessagesProxyI
	c  ConnectionEventsI
	e  EventsI
	se SendStrategyI
	i  RoomInformationI
}

// Init configure dependencies with other components of the room
func (room *RoomAPI) Init(builder ComponentBuilderI) {
	builder.BuildSync(&room.s)
	builder.BuildMessages(&room.m)
	builder.BuildRoomConnectionEvents(&room.c)
	builder.BuildEvents(&room.e)
	builder.BuildSender(&room.se)
	builder.BuildInformation(&room.i)
}

// Handle processes the request came from the user
func (room *RoomAPI) Handle(conn *Connection, rr *RoomRequest) {
	go room.s.do(func() {
		if rr.IsGet() {
			room.GetRoom(conn)
		} else if rr.IsSend() {
			room.handleSent(conn, rr.Send)
		} else if rr.Message != nil {
			room.handleMessage(conn, rr.Message)
		}
	})
}

func (room *RoomAPI) handleMessage(conn *Connection, message *models.Message) {
	room.s.doWithConn(conn, func() {
		/*
			Message(room.i.lobby, conn, message, room.m.appendMessage,
				room.m.setMessage, room.m.removeMessage, room.m.findMessage,
				room.r.send.sendAll, room.r.All, room.r, room.m.dbChatID)*/
		HandleMessage(conn, message, room.m)
	})
}

func (room *RoomAPI) handleSent(conn *Connection, request *RoomSend) {
	switch {
	case request.Messages != nil:
		room.GetMessages(conn, request.Messages)
	case request.Cell != nil:
		room.PostCell(conn, request.Cell)
	case request.Action != nil:
		room.PostAction(conn, *request.Action)
	}
}

// GetRoom handle "GET /room", return all room information
func (room *RoomAPI) GetRoom(conn *Connection) {
	room.s.doWithConn(conn, func() {
		room.se.Room(conn)
	})
}

// GetMessages handle "GET /messages", return all room messages
func (room *RoomAPI) GetMessages(conn *Connection, settings *models.Messages) {
	room.s.doWithConn(conn, func() {
		Messages(conn, settings, room.m.Messages())
	})
}

// PostCell  handle "POST /cell" processes the Cell came from the user
func (room *RoomAPI) PostCell(conn *Connection, cell *Cell) {
	room.s.doWithConn(conn, func() {
		room.e.OpenCell(conn, cell)
	})
}

// PostAction handle "POST /action" processes the Cell came from the user
func (room *RoomAPI) PostAction(conn *Connection, action int) {
	room.s.doWithConn(conn, func() {
		switch action {
		case ActionBackToLobby:
			room.c.Leave(conn)
		case ActionDisconnect:
			room.c.Disconnect(conn)
		case ActionReconnect:
			room.c.Reconnect(conn)
		case ActionGiveUp:
			room.c.GiveUp(conn)
		case ActionRestart:
			room.c.Restart(conn)
		}
	})
}
