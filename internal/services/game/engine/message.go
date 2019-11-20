package engine

import (
	"math/rand"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/handlers"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/clients"
	cmodels "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/database"
	chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/game/metrics"
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
func GetChatID(chatS clients.Chat, chatType int32, typeID int32) (int32, error) {
	var (
		newChat = &chat.Chat{
			Type:   chatType,
			TypeId: typeID,
		}
		chatID *chat.ChatID
		err    error
	)

	chatID, err = chatS.GetChat(newChat)

	if err != nil {
		utils.Debug(false, "cant access to chat service", err.Error())
		return 0, err
	}
	return chatID.Value, err
}

// GetChatIDAndMessages accesses the chat service to get the ID of the chat and
// all messages
func GetChatIDAndMessages(chatS clients.Chat, loc *time.Location, chatType, typeID int32,
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

	chatID, err = chatS.GetChat(newChat)

	if err != nil {
		utils.Debug(false, "cant access to chat service", err.Error())
		return 0, nil, err
	}
	pMessages, err = chatS.GetChatMessages(chatID)
	if err != nil {
		utils.Debug(false, "cant get messages!", err.Error())
		return 0, nil, err
	}

	var messages []*models.Message
	messages, err = cmodels.MessagesFromProto(loc, pMessages.Messages...)

	for _, message := range messages {
		setImage(message.User)
	}

	return chatID.Value, messages, err
}

func HandleMessage(conn *Connection,
	message *models.Message, handler MessagesI) error {
	handler.Fix(message, conn)
	msg, err := handler.Proto(message)
	if err != nil {
		return err
	}

	// ignore models.StartWrite, models.FinishWrite
	switch message.Action {
	case models.Write:
		err = handler.Write(message, msg)
	case models.Update:
		err = handler.Update(message, msg)
	case models.Delete:
		err = handler.Delete(message, msg)
	}
	if err != nil {
		handler.HandleError(message, msg)
	} else {
		handler.Send(message)
	}
	return err
}

// Message send message to connections
func Message(chatS clients.Chat, lobby *Lobby, conn *Connection, message *models.Message,
	append AppendMessage, update UpdateMessage, delete DeleteMessage,
	find FindMessage, send Sender, predicate SendPredicate, room *Room,
	chatID int32) (err error) {

	message.User = conn.User

	message.Time = time.Now().In(lobby.location())

	if message.Action == models.Write {
		rand.Seed(time.Now().UnixNano())
		message.ID = rand.Int31n(10000000) // в конфиг?
	}

	msg, err := cmodels.MessageToProto(message, chatID)

	if err != nil {
		return err
	}

	var msgID *chat.MessageID

	// ignore models.StartWrite, models.FinishWrite
	switch message.Action {
	case models.Write:
		append(message)
		msgID, err = chatS.AppendMessage(msg)
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

		_, err = chatS.UpdateMessage(msg)

	case models.Delete:
		if message.ID <= 0 {
			return re.InvalidMessageID()
		}
		delete(find(message))

		_, err = chatS.DeleteMessage(msg)

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
					if room.messages.ChatID() != 0 {
						return room.messages.ChatID(), nil
					}
					id, err := GetChatID(chatS, cmodels.RoomType, room.info.RoomID())
					if err != nil {
						room.messages.setChatID(id)
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
	messages.Messages = messageSlice[messages.Offset : messages.Offset+messages.Limit]
	messages.Capacity = size

	response := models.Response{
		Type:  "GameMessages",
		Value: messages,
	}

	conn.SendInformation(&response)
}
