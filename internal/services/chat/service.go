package chat

import (
	"sync"
	"time"

	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"context"

	"database/sql"

	//
	_ "github.com/lib/pq"
)

// Service of chat
type Service struct {
	wGroup *sync.WaitGroup
	ID     int
	DB     *sql.DB
}

//NewService Create new instance of service
func NewService(db *sql.DB, id int) *Service {
	return &Service{
		wGroup: &sync.WaitGroup{},
		DB:     db,
		ID:     id,
	}
}

func (service *Service) Debug(needPanic bool, text ...interface{}) {
	utils.Debug(needPanic, "Chat", service.ID, ":", text)
}

func (service *Service) Check() (bool, error) {
	err := service.DB.Ping()
	if err != nil {
		return false, err
	}

	return false, nil
}

func (service *Service) Close() {
	timeout := 2 * time.Second //TODO в конфиг!
	utils.WaitWithTimeout(service.wGroup, timeout)

	service.DB.Close()
	return
}

// CreateChat create chat with or without users.
// Specify the type of chat and id received from the corresponding database table
// Return id for this chat, save it. It must be transferred to any chat operations
func (service *Service) CreateChat(ctx context.Context, chat *ChatWithUsers) (*ChatID, error) {
	service.wGroup.Add(1)
	defer func() {
		service.wGroup.Done()
	}()

	if chat == nil {
		return &ChatID{}, re.InvalidMessage()
	}

	var (
		tx     *sql.Tx
		err    error
		ChatID = &ChatID{}
	)

	if tx, err = service.DB.Begin(); err != nil {
		return ChatID, err
	}
	defer tx.Rollback()

	if ChatID, err = service.insertChat(tx, chat.Type, chat.TypeId); err != nil {
		return ChatID, err
	}

	if err = service.insertUsers(tx, ChatID.Value, chat.Users...); err != nil {
		return ChatID, err
	}

	if err = tx.Commit(); err != nil {
		return ChatID, err
	}
	service.Debug(false, "CreateChat success")
	return ChatID, err
}

// GetChat get the ID of the chat, based on its type and the passed ID of this type
func (service *Service) GetChat(ctx context.Context, chat *Chat) (*ChatID, error) {
	service.wGroup.Add(1)
	defer func() {
		service.wGroup.Done()
	}()

	if chat == nil {
		return &ChatID{}, re.InvalidChatID()
	}

	id, err := service.getChat(chat)
	if err == nil {
		service.Debug(false, "GetChat success:", chat.Type, chat.TypeId, id)
	} else {
		var (
			tx *sql.Tx
		)
		if tx, err = service.DB.Begin(); err != nil {
			return id, err
		}
		defer tx.Rollback()
		if id, err = service.insertChat(tx, chat.Type, chat.TypeId); err != nil {
			return id, err
		}
		if err = tx.Commit(); err != nil {
			return id, err
		}
	}

	return id, err
}

// AppendMessage append message to database
// to work correctly, specify the ID of the chat(in the message) in which
// the operation occurs
// Return id for this message, save it. It must be transferred to any message
// operations
func (service *Service) AppendMessage(ctx context.Context, message *Message) (*MessageID, error) {
	service.wGroup.Add(1)
	defer func() {
		service.wGroup.Done()
	}()

	if message == nil {
		return &MessageID{}, re.InvalidMessage()
	}

	if message.ChatId <= 0 {
		return &MessageID{}, re.InvalidMessageID()
	}

	id, err := service.insertMessage(message)
	if err == nil {
		service.Debug(false, "AppendMessage success")
	}
	return id, err
}

// AppendMessages append messages to database
func (service *Service) AppendMessages(ctx context.Context, messages *Messages) (*MessagesID, error) {
	service.wGroup.Add(1)
	defer func() {
		service.wGroup.Done()
	}()

	if messages == nil {
		return &MessagesID{}, re.InvalidMessage()
	}

	ids, err := service.insertMessages(messages)
	if err == nil {
		service.Debug(false, "AppendMessages success")
	}
	return ids, err
}

// UpdateMessage update message in database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (service *Service) UpdateMessage(ctx context.Context, message *Message) (*Result, error) {
	service.wGroup.Add(1)
	defer func() {
		service.wGroup.Done()
	}()

	if message == nil {
		return &Result{}, re.InvalidMessage()
	}

	if message.Id <= 0 {
		return &Result{}, re.InvalidMessageID()
	}

	res, err := service.updateMessage(message)
	if err == nil {
		service.Debug(false, "UpdateMessage success")
	}
	return res, err
}

// DeleteMessage delete message from database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (service *Service) DeleteMessage(ctx context.Context, message *Message) (*Result, error) {
	service.wGroup.Add(1)
	defer func() {
		service.wGroup.Done()
	}()

	if message == nil {
		return &Result{}, re.InvalidMessage()
	}

	if message.Id <= 0 {
		return &Result{}, re.InvalidMessageID()
	}

	res, err := service.deleteMessage(message)
	if err == nil {
		service.Debug(false, "deleteMessage success")
	}
	return res, err
}

// InviteToChat invite user to the chat
// to work correctly, specify user and id of the chat
func (service *Service) InviteToChat(ctx context.Context, userInChat *UserInGroup) (*Result, error) {
	service.wGroup.Add(1)
	defer func() {
		service.wGroup.Done()
	}()

	if userInChat == nil {
		return &Result{}, re.InvalidUser()
	}

	if userInChat.User == nil || userInChat.User.Id <= 0 {
		return &Result{}, re.InvalidUser()
	}

	if userInChat.Chat == nil || userInChat.Chat.Id <= 0 {
		return &Result{}, re.InvalidChatID()
	}

	var (
		tx  *sql.Tx
		err error
	)

	if tx, err = service.DB.Begin(); err != nil {
		return &Result{Done: false}, err
	}
	defer tx.Rollback()

	if err = service.insertUsers(tx, userInChat.Chat.Id, userInChat.User); err != nil {
		return &Result{Done: false}, err
	}

	err = tx.Commit()
	if err == nil {
		service.Debug(false, "InviteToChat success")
	}

	return &Result{Done: true}, err
}

// LeaveChat leave user from the chat
// to work correctly, specify user and id of the chat
func (service *Service) LeaveChat(ctx context.Context, userInChat *UserInGroup) (*Result, error) {
	service.wGroup.Add(1)
	defer func() {
		service.wGroup.Done()
	}()

	if userInChat == nil {
		return &Result{}, re.InvalidUser()
	}

	if userInChat.User == nil || userInChat.User.Id <= 0 {
		return &Result{}, re.InvalidUser()
	}

	if userInChat.Chat == nil || userInChat.Chat.Id <= 0 {
		return &Result{}, re.InvalidChatID()
	}

	res, err := service.deleteUserInChat(userInChat)
	if err == nil {
		service.Debug(false, "LeaveChat success")
	}
	return res, err
}

// GetChatMessages get all messages from the chad with specified id
func (service Service) GetChatMessages(ctx context.Context, chatID *ChatID) (*Messages, error) {
	service.wGroup.Add(1)
	defer func() {
		service.wGroup.Done()
	}()

	if chatID == nil {
		return &Messages{}, re.InvalidMessage()
	}

	messages, err := service.getChatMessages(chatID)
	if err == nil {
		service.Debug(false, "GetChatMessages success. Id was", chatID.Value,
			"Messages amount:", len(messages.Messages))
	}
	return messages, err
}
