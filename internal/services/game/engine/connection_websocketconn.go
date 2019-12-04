package engine

import (
	"bufio"
	"io"
	"time"
)

// WebsocketConnI the interface of the websocket connection
//	 Is based on the websocket.Conn
type WebsocketConnI interface {
	WriteMessage(messageType int, data []byte) error
	NextWriter(messageType int) (io.WriteCloser, error)
	SetWriteDeadline(t time.Time) error

	ReadMessage() (messageType int, p []byte, err error)
	SetReadLimit(limit int64)
	SetReadDeadline(t time.Time) error

	SetPongHandler(h func(appData string) error)

	Close() error
}

// WriteCloserStub stu
type WriteCloserStub struct {
	*bufio.Writer
}

// Close stub
func (mwc *WriteCloserStub) Close() error {
	return nil
}

// WebsocketConnStub stub of websocket.Conn
type WebsocketConnStub struct {
}

// WriteMessage stub
func (WebsocketConnStub *WebsocketConnStub) WriteMessage(messageType int, data []byte) error {
	return nil
}

// NextWriter stub
func (WebsocketConnStub *WebsocketConnStub) NextWriter(messageType int) (io.WriteCloser, error) {
	return &WriteCloserStub{}, nil
}

// SetWriteDeadline stub
func (WebsocketConnStub *WebsocketConnStub) SetWriteDeadline(t time.Time) error {
	return nil
}

// ReadMessage stub
func (WebsocketConnStub *WebsocketConnStub) ReadMessage() (messageType int, p []byte, err error) {
	return 0, []byte{}, nil
}

// SetReadLimit stub
func (WebsocketConnStub *WebsocketConnStub) SetReadLimit(limit int64) {}

// SetReadDeadline stub
func (WebsocketConnStub *WebsocketConnStub) SetReadDeadline(t time.Time) error {
	return nil
}

// SetPongHandler stub
func (WebsocketConnStub *WebsocketConnStub) SetPongHandler(h func(appData string) error) {

}

// Close stub
func (WebsocketConnStub *WebsocketConnStub) Close() error {
	return nil
}
