package game

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

// Sender is the func that send information to connections
type Sender func(interface{}, SendPredicate)

type AppendMessage func(*models.Message)
type UpdateMessage func(int, *models.Message)
type FindMessage func(*models.Message) int
type DeleteMessage func(int)

// Message send message to connections
func Message(lobby *Lobby, conn *Connection, message *models.Message,
	append AppendMessage, update UpdateMessage, delete DeleteMessage,
	find FindMessage, send Sender, predicate SendPredicate, inRoom bool,
	roomID string) (err error) {
	message.User = conn.User
	message.Time = time.Now()

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
