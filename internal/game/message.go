package game

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

// Sender is the func that send information to connections
type Sender func(interface{}, SendPredicate)

// AppendMessage - the function to append message to message slice
type AppendMessage func(*models.Message)

// UpdateMessage - the function to update message in message slice
type UpdateMessage func(int, *models.Message)

// FindMessage - the function to search message in message slice
type FindMessage func(*models.Message) int

// DeleteMessage - the function to delete message from message slice
type DeleteMessage func(int)

// Message send message to connections
func Message(lobby *Lobby, conn *Connection, message *models.Message,
	append AppendMessage, update UpdateMessage, delete DeleteMessage,
	find FindMessage, send Sender, predicate SendPredicate, inRoom bool,
	roomID string) (err error) {
	message.User = conn.User

	loc, _ := time.LoadLocation("Europe/Moscow")
	message.Time = time.Now().In(loc)

	// ignore models.StartWrite, models.FinishWrite
	switch message.Action {
	case models.Write:
		append(message)
		message.ID, err = lobby.db.CreateMessage(message, inRoom, roomID)
	case models.Update:
		if message.ID <= 0 {
			return re.ErrorMessageInvalidID()
		}
		update(find(message), message)
		_, err = lobby.db.UpdateMessage(message)
	case models.Delete:
		if message.ID <= 0 {
			return re.ErrorMessageInvalidID()
		}
		delete(find(message))
		_, err = lobby.db.DeleteMessage(message)
	}
	if err != nil {
		return err
	}

	response := models.Response{
		Type:  "GameMessage",
		Value: message,
	}

	send(response, predicate)
	return err

}

// Messages processes the receipt of an object Messages from the user
func Messages(conn *Connection, messages *models.Messages,
	messageSlice []*models.Message) {

	size := len(messageSlice)
	if messages.Offset < 0 || messages.Offset >= size {
		messages.Offset = 0
	}
	if messages.Limit < 0 || messages.Limit > size {
		messages.Limit = size
	}
	messages.Messages = messageSlice[messages.Offset:messages.Limit]
	messages.Capacity = size

	response := models.Response{
		Type:  "GameMessages",
		Value: messages,
	}

	conn.SendInformation(response)
}
