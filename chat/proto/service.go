package chat

import (
	"context"

	"database/sql"

	//
	_ "github.com/lib/pq"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

type Service struct {
	DB *sql.DB
}

func (service *Service) CreateChat(ctx context.Context, chat *ChatWithUsers) (*ChatID, error) {
	var (
		tx     *sql.Tx
		err    error
		ChatID = &ChatID{}
	)

	if tx, err = service.DB.Begin(); err != nil {
		return ChatID, err
	}
	defer tx.Rollback()

	if ChatID, err = service.insertChat(tx, chat); err != nil {
		return ChatID, err
	}

	if err = service.insertUsers(tx, ChatID.Value, chat.Users...); err != nil {
		return ChatID, err
	}

	if err = tx.Commit(); err != nil {
		return ChatID, err
	}
	utils.Debug(false, "CreateChat success")
	return ChatID, err
}

func (service *Service) GetChat(ctx context.Context, chat *Chat) (*ChatID, error) {
	id, err := service.getChat(chat)
	if err == nil {
		utils.Debug(false, "GetChat success:", chat.Type, chat.TypeId, id)
	}
	return id, err
}

func (service *Service) AppendMessage(ctx context.Context, message *Message) (*MessageID, error) {
	utils.Debug(false, "Message time", message.Time)
	id, err := service.insertMessage(message)
	if err == nil {
		utils.Debug(false, "AppendMessage success")
	}
	return id, err
}
func (service *Service) AppendMessages(ctx context.Context, messages *Messages) (*MessagesID, error) {
	ids, err := service.insertMessages(messages)
	if err == nil {
		utils.Debug(false, "AppendMessages success")
	}
	return ids, err
}
func (service *Service) UpdateMessage(ctx context.Context, message *Message) (*Result, error) {
	res, err := service.updateMessage(message)
	if err == nil {
		utils.Debug(false, "UpdateMessage success")
	}
	return res, err
}
func (service *Service) DeleteMessage(ctx context.Context, message *Message) (*Result, error) {
	res, err := service.deleteMessage(message)
	if err == nil {
		utils.Debug(false, "deleteMessage success")
	}
	return res, err
}

func (service *Service) InviteToChat(ctx context.Context, userInChat *UserInGroup) (*Result, error) {
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
		utils.Debug(false, "InviteToChat success")
	}

	return &Result{Done: true}, err
}
func (service *Service) LeaveChat(ctx context.Context, userInChat *UserInGroup) (*Result, error) {
	res, err := service.deleteUserInChat(userInChat)
	if err == nil {
		utils.Debug(false, "LeaveChat success")
	}
	return res, err
}

func (service Service) GetChatMessages(ctx context.Context, chatID *ChatID) (*Messages, error) {
	messages, err := service.getChatMessages(chatID)
	if err == nil {
		utils.Debug(false, "GetChatMessages success. Id was", chatID.Value,
			"Messages amount:", len(messages.Messages))
	}
	return messages, err
}
