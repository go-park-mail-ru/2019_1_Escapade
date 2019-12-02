package engine

import "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"

// codes of connection event
const (
	Join      = iota
	EnterRoom = iota
	Leave     = iota
	Action    = iota
)

type ConnEventMsg struct {
	Code    int
	Content interface{}
}

type ConnEvents struct {
	eventsChan chan ConnEventMsg
}

func NewConnEvent() *ConnEvents {
	return &ConnEvents{make(chan ConnEventMsg)}
}

func (conn *ConnEvents) Close() {
	close(conn.eventsChan)
}

func (conn *ConnEvents) Get() ConnEventMsg {
	return <-conn.eventsChan
}

func (conn *ConnEvents) Join() {
	conn.eventsChan <- ConnEventMsg{
		Code:    Join,
		Content: nil,
	}
}

func (conn *ConnEvents) EnterRoom(roomID string) {
	rs := &models.RoomSettings{}
	rs.ID = roomID
	conn.eventsChan <- ConnEventMsg{
		Code:    EnterRoom,
		Content: rs,
	}
}

func (conn *ConnEvents) Leave() {
	conn.eventsChan <- ConnEventMsg{
		Code:    Leave,
		Content: nil,
	}
}

func (conn *ConnEvents) Request(content []byte) {
	conn.eventsChan <- ConnEventMsg{
		Code:    Action,
		Content: content,
	}
}
