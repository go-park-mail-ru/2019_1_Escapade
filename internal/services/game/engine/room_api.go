package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

type RoomAPI struct {
	s  SyncI
	m  *RoomMessages
	ce *RoomConnectionEvents
	e  *RoomEvents
	se *RoomSender
	i  *RoomInformation
}

func (room *RoomAPI) Init(s SyncI,
	m *RoomMessages, ce *RoomConnectionEvents,
	se *RoomSender, e *RoomEvents, i *RoomInformation) {
	room.s = s
	room.m = m
	room.ce = ce
	room.se = se
	room.e = e
	room.i = i
}

// HandleRequest processes the request came from the user
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

func (room *RoomAPI) GetRoom(conn *Connection) {
	room.s.doWithConn(conn, func() {
		room.se.Room(conn)
	})
}

func (room *RoomAPI) GetMessages(conn *Connection, settings *models.Messages) {
	room.s.doWithConn(conn, func() {
		Messages(conn, settings, room.m.Messages())
	})
}

// CellHandle processes the Cell came from the user
func (room *RoomAPI) PostCell(conn *Connection, cell *Cell) {
	room.s.doWithConn(conn, func() {
		room.e.OpenCell(conn, cell)
	})
}

func (room *RoomAPI) PostAction(conn *Connection, action int) {
	room.s.doWithConn(conn, func() {
		switch action {
		case ActionBackToLobby:
			room.ce.Leave(conn)
		case ActionDisconnect:
			room.ce.Disconnect(conn)
		case ActionReconnect:
			room.ce.Reconnect(conn)
		case ActionGiveUp:
			room.ce.GiveUp(conn)
		case ActionRestart:
			room.ce.Restart(conn)
		}
	})
}
