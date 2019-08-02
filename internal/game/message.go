package game

import (
	"context"
	"fmt"
	"time"

	сhat "github.com/go-park-mail-ru/2019_1_Escapade/chat/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// Sender is the func that send information to connections
type Sender func(utils.JSONtype, SendPredicate)

// AppendMessage - the function to append message to message slice
type AppendMessage func(*models.Message)

// UpdateMessage - the function to update message in message slice
type UpdateMessage func(int, *models.Message)

// FindMessage - the function to search message in message slice
type FindMessage func(*models.Message) int

// DeleteMessage - the function to delete message from message slice
type DeleteMessage func(int)

func GetChatIDAndMessages(loc *time.Location, chatType сhat.ChatType, typeID int32) (int32, []*models.Message, error) {
	var (
		newChat = &сhat.Chat{
			Type:   chatType,
			TypeId: typeID,
		}
		chatID    *сhat.ChatID
		pMessages *сhat.Messages
		err       error
	)

	chatID, err = clients.ALL.Chat.GetChat(context.Background(), newChat)

	if err != nil {
		utils.Debug(true, "cant access to chat service", err.Error())
		return 0, nil, err
	}
	pMessages, err = clients.ALL.Chat.GetChatMessages(context.Background(), chatID)
	if err != nil {
		utils.Debug(true, "cant get messages!", err.Error())
		return 0, nil, err
	}

	var messages []*models.Message
	messages = сhat.MessagesFromProto(loc, pMessages.Messages...)
	//db.getMessages(tx, true, game.RoomID)

	return chatID.Value, messages, err
}

// Message send message to connections
func Message(lobby *Lobby, conn *Connection, message *models.Message,
	append AppendMessage, update UpdateMessage, delete DeleteMessage,
	find FindMessage, send Sender, predicate SendPredicate, inRoom bool, chatID int32) (err error) {
	message.User = conn.User

	message.Time = time.Now().In(lobby.location())

	// ignore models.StartWrite, models.FinishWrite
	switch message.Action {
	case models.Write:
		append(message)
		utils.Debug(false, "look at time", message.Time)
		msg := сhat.MessageToProto(message)
		utils.Debug(false, "newtime", msg.Time)
		msg.ChatId = chatID
		msgID, err := clients.ALL.Chat.AppendMessage(context.Background(), msg)
		if err != nil {
			return err
		}
		message.ID = msgID.Value
		//message.ID, err = lobby.db.CreateMessage(message, inRoom, roomID)
		if lobby.config().Metrics {
			if inRoom {
				metrics.RoomsMessages.Inc()
			} else {
				metrics.LobbyMessages.Inc()
			}
		}
	case models.Update:
		if message.ID <= 0 {
			return re.ErrorMessageInvalidID()
		}
		update(find(message), message)

		msg := сhat.MessageToProto(message)
		_, err := clients.ALL.Chat.UpdateMessage(context.Background(), msg)
		if err != nil {
			return err
		}

		//_, err = lobby.db.UpdateMessage(message)
	case models.Delete:
		if message.ID <= 0 {
			return re.ErrorMessageInvalidID()
		}
		delete(find(message))

		msg := сhat.MessageToProto(message)
		_, err := clients.ALL.Chat.DeleteMessage(context.Background(), msg)
		if err != nil {
			return err
		}

		//_, err = lobby.db.DeleteMessage(message)
		if lobby.config().Metrics {
			if inRoom {
				metrics.RoomsMessages.Dec()
			} else {
				metrics.LobbyMessages.Dec()
			}
		}
	}
	if err != nil {
		return err
	}

	response := models.Response{
		Type:  "GameMessage",
		Value: message,
	}

	fmt.Println("response")

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
