package engine

import (
	"context"
	"math/rand"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/clients"
	handlers "github.com/go-park-mail-ru/2019_1_Escapade/internal/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// Sender is the func that send information to connections
type Sender func(handlers.JSONtype, SendPredicate)

// AppendMessage - the function to append message to message slice
type AppendMessage func(*models.Message)

// UpdateMessage - the function to update message in message slice
type UpdateMessage func(int, *models.Message)

// FindMessage - the function to search message in message slice
type FindMessage func(*models.Message) int

// DeleteMessage - the function to delete message from message slice
type DeleteMessage func(int)

// GetChatID accesses the chat service to get the ID of the chat
func GetChatID(chatType chat.ChatType, typeID int32) (int32, error) {
	var (
		newChat = &chat.Chat{
			Type:   chatType,
			TypeId: typeID,
		}
		chatID *chat.ChatID
		err    error
	)

	chatID, err = clients.ALL.Chat().GetChat(context.Background(), newChat)

	if err != nil {
		utils.Debug(false, "cant access to chat service", err.Error())
		return 0, err
	}
	return chatID.Value, err
}

// GetChatIDAndMessages accesses the chat service to get the ID of the chat and
// all messages
func GetChatIDAndMessages(loc *time.Location, chatType chat.ChatType, typeID int32,
	setImage SetImage) (int32, []*models.Message, error) {
	var (
		newChat = &chat.Chat{
			Type:   chatType,
			TypeId: typeID,
		}
		chatID    *chat.ChatID
		pMessages *chat.Messages
		err       error
	)

	chatID, err = clients.ALL.Chat().GetChat(context.Background(), newChat)

	if err != nil {
		utils.Debug(false, "cant access to chat service", err.Error())
		return 0, nil, err
	}
	pMessages, err = clients.ALL.Chat().GetChatMessages(context.Background(), chatID)
	if err != nil {
		utils.Debug(false, "cant get messages!", err.Error())
		return 0, nil, err
	}

	var messages []*models.Message
	messages, err = chat.MessagesFromProto(loc, pMessages.Messages...)

	for _, message := range messages {
		setImage(message.User)
	}

	return chatID.Value, messages, err
}

// Message send message to connections
func Message(lobby *Lobby, conn *Connection, message *models.Message,
	append AppendMessage, update UpdateMessage, delete DeleteMessage,
	find FindMessage, send Sender, predicate SendPredicate, room *Room,
	chatID int32) (err error) {

	message.User = conn.User

	message.Time = time.Now().In(lobby.location())

	if message.Action == models.Write {
		rand.Seed(time.Now().UnixNano())
		message.ID = rand.Int31n(10000000) // в конфиг?
	}

	msg, err := chat.MessageToProto(message)

	if err != nil {
		return err
	}
	msg.ChatId = chatID

	var msgID *chat.MessageID

	// ignore models.StartWrite, models.FinishWrite
	switch message.Action {
	case models.Write:
		append(message)
		msgID, err = clients.ALL.Chat().AppendMessage(context.Background(), msg)
		if msgID != nil {
			message.ID = msgID.Value
		}
		if lobby.config().Metrics {
			if room != nil {
				metrics.RoomsMessages.Inc()
			} else {
				metrics.LobbyMessages.Inc()
			}
		}
	case models.Update:
		if message.ID <= 0 {
			return re.InvalidMessageID()
		}
		update(find(message), message)

		_, err = clients.ALL.Chat().UpdateMessage(context.Background(), msg)

	case models.Delete:
		if message.ID <= 0 {
			return re.InvalidMessageID()
		}
		delete(find(message))

		_, err = clients.ALL.Chat().DeleteMessage(context.Background(), msg)

		if lobby.config().Metrics {
			if room != nil {
				metrics.RoomsMessages.Dec()
			} else {
				metrics.LobbyMessages.Dec()
			}
		}
	}
	if err != nil {
		action := message.Action
		if room != nil {
			lobby.AddNotSavedMessage(&MessageWithAction{
				message, msg, action, func() (int32, error) {
					if room.dbChatID != 0 {
						return room.dbChatID, nil
					}
					id, err := GetChatID(chat.ChatType_ROOM, room.dbRoomID)
					if err != nil {
						room.dbChatID = id
					}
					return id, err
				}})
		} else {
			lobby.AddNotSavedMessage(&MessageWithAction{
				message, msg, action, func() (int32, error) {
					dbChatID := lobby.dbChatID()
					if dbChatID != 0 {
						return dbChatID, nil
					}
					return dbChatID, re.InvalidChatID()
				}})
		}
	}

	sendMessages(send, predicate, message)
	return err
}

func sendMessages(send Sender, predicate SendPredicate, messages ...*models.Message) {
	for _, message := range messages {
		response := models.Response{
			Type:  "GameMessage",
			Value: message,
		}

		send(&response, predicate)
	}
}

func sendMessagesTodelete(send Sender, predicate SendPredicate, messages ...*models.Message) {
	for _, message := range messages {
		newMessage := models.Message{
			ID:     message.ID,
			Action: models.Delete,
		}
		response := models.Response{
			Type:  "GameMessage",
			Value: newMessage,
		}

		send(&response, predicate)
	}
}

// Messages processes the receipt of an object Messages from the user
// TODO: delete or normally implement
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

	conn.SendInformation(&response)
}
