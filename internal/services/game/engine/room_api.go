package engine

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

type RoomAPI struct {
	r *Room
	s SyncI
	m *RoomMessages
}

func (room *RoomAPI) Init(r *Room, s SyncI, m *RoomMessages) {
	room.r = r
	room.s = s
	room.m = m
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
		if conn.Index() < 0 {
			message.Status = models.StatusObserver
		} else {
			message.Status = models.StatusPlayer
		}
		Message(room.r.lobby, conn, message, room.m.appendMessage,
			room.m.setMessage, room.m.removeMessage, room.m.findMessage,
			room.r.send.sendAll, room.r.All, room.r, room.m.dbChatID)
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
		room.r.send.Room(conn)
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
		if !room.r.people.isAlive(conn) {
			return
		}
		status := room.r.events.Status()
		if status == StatusFlagPlacing {
			room.r.field.SetFlag(conn, cell)
		} else if status == StatusRunning {
			room.r.field.OpenCell(conn, cell)
		}
	})
}

func (room *RoomAPI) PostAction(conn *Connection, action int) {
	room.s.doWithConn(conn, func() {
		switch action {
		case ActionBackToLobby:
			room.r.connEvents.Leave(conn)
		case ActionDisconnect:
			room.r.connEvents.Disconnect(conn)
		case ActionReconnect:
			room.r.connEvents.Reconnect(conn)
		case ActionGiveUp:
			room.r.connEvents.GiveUp(conn)
		case ActionRestart:
			room.r.connEvents.Restart(conn)
		}
	})
}
