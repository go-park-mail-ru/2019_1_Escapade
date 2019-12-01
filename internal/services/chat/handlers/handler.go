package handlers

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"

	cdb "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

// repositories stores all implementations of operations in the database
type repositories struct {
	user    cdb.UserRepositoryI
	chat    cdb.ChatRepositoryI
	message cdb.MessageRepositoryI
}

type Handler struct {
	user    cdb.UserUseCaseI
	chat    cdb.ChatUseCaseI
	message cdb.MessageUseCaseI
}

// InitWithPostgreSQL apply postgreSQL as database
func (h *Handler) InitWithPostgreSQL(c *config.Configuration) error {
	var (
		reps = repositories{
			user:    &cdb.UserRepositoryPQ{},
			message: &cdb.MessageRepositoryPQ{},
			chat:    &cdb.ChatRepositoryPQ{},
		}
		database = &database.PostgresSQL{}
	)
	return h.Init(c, database, reps)
}

func (h *Handler) Init(c *config.Configuration, db database.DatabaseI, reps repositories) error {
	err := db.Open(c.DataBase)
	if err != nil {
		return err
	}

	var user = &cdb.UserUseCase{}
	user.Init(reps.user)
	h.user = user
	err = h.user.Use(db)
	if err != nil {
		return err
	}

	var message = &cdb.MessageUseCase{}
	message.Init(reps.message)
	h.message = message
	err = h.message.Use(db)
	if err != nil {
		return err
	}

	var chat = &cdb.ChatUseCase{}
	chat.Init(reps.chat, reps.user)
	h.chat = chat
	err = h.chat.Use(db)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) Check() (bool, error) {
	err := h.user.Get().Ping()
	if err != nil {
		return false, err
	}

	return false, nil
}

func (h *Handler) Close() {
	h.user.Close()
	h.message.Close()
	h.chat.Close()
	return
}

// CreateChat chat with or without users.
// Specify the type of chat and id received from the corresponding database
// table. Return id for this chat, save it. It must be transferred to any
// chat operations
func (h *Handler) CreateChat(ctx context.Context,
	in *proto.ChatWithUsers) (*proto.ChatID, error) {
	return h.chat.Create(in)
}

// GetChat get the ID of the chat, based on its type and the passed ID of
// this type
func (h *Handler) GetChat(ctx context.Context,
	in *proto.Chat) (*proto.ChatID, error) {
	return h.chat.GetOne(in)
}

// InviteToChat invite user to the chat
// to work correctly, specify user and id of the chat
func (h *Handler) InviteToChat(ctx context.Context,
	in *proto.UserInGroup) (*proto.Result, error) {
	return h.user.InviteToChat(in)
}

// LeaveChat leave user from the chat
// to work correctly, specify user and id of the chat
func (h *Handler) LeaveChat(ctx context.Context,
	in *proto.UserInGroup) (*proto.Result, error) {
	return h.user.LeaveChat(in)
}

// AppendMessage to database
// to work correctly, specify the ID of the chat(in the message) in which
// the operation occurs
// Return id for this message, save it. It must be transferred to any message
// operations
func (h *Handler) AppendMessage(ctx context.Context,
	in *proto.Message) (*proto.MessageID, error) {
	return h.message.AppendOne(in)
}

// AppendMessages to database
func (h *Handler) AppendMessages(ctx context.Context,
	in *proto.Messages) (*proto.MessagesID, error) {
	return h.message.AppendMany(in)
}

// UpdateMessage in database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (h *Handler) UpdateMessage(ctx context.Context,
	in *proto.Message) (*proto.Result, error) {
	return h.message.Update(in)
}

// DeleteMessage from database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (h *Handler) DeleteMessage(ctx context.Context,
	in *proto.Message) (*proto.Result, error) {
	return h.message.Delete(in)
}

// GetChatMessages get all messages from the chad with specified id
func (h *Handler) GetChatMessages(ctx context.Context,
	in *proto.ChatID) (*proto.Messages, error) {
	return h.message.GetAll(in)
}
