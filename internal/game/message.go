package game

import (
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// Sender is the func that send information to connections
type Sender func(interface{}, SendPredicate)

type SetToMessages func(message *models.Message)

// Message send message to connections
func Message(lobby *Lobby, conn *Connection,
	message *models.Message, setToMessages SetToMessages,
	send Sender, predicate SendPredicate, inRoom bool, roomID string) {
	message.User = conn.User
	message.Time = time.Now()
	setToMessages(message)
	lobby.db.SaveMessage(message, inRoom, roomID)

	response := models.Response{
		Type:  "GameMessage",
		Value: message,
	}
	send(response, predicate)

}
