package handlers

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

type Handler struct {
	db *database.Input
}

// Init initialize Handler
func (h *Handler) Init(c config.Database, input *database.Input) error {

	input.Init()
	if err := input.IsValid(); err != nil {
		return err
	}

	h.db = input
	return h.db.Connect(c)
}

// Check health of service
func (h *Handler) Check() (bool, error) {
	return false, h.db.UserUC.Get().Ping()
}

// Close connection to database
func (h *Handler) Close() error {
	return h.db.Close()
}

// CreateChat chat with or without users.
// Specify the type of chat and id received from the corresponding database
// table. Return id for this chat, save it. It must be transferred to any
// chat operations
func (h *Handler) CreateChat(ctx context.Context,
	in *proto.ChatWithUsers) (*proto.ChatID, error) {
	return h.db.ChatUC.Create(in)
}

// GetChat get the ID of the chat, based on its type and the passed ID of
// this type
func (h *Handler) GetChat(ctx context.Context,
	in *proto.Chat) (*proto.ChatID, error) {
	return h.db.ChatUC.GetOne(in)
}

// InviteToChat invite user to the chat
// to work correctly, specify user and id of the chat
func (h *Handler) InviteToChat(ctx context.Context,
	in *proto.UserInGroup) (*proto.Result, error) {
	return h.db.UserUC.InviteToChat(in)
}

// LeaveChat leave user from the chat
// to work correctly, specify user and id of the chat
func (h *Handler) LeaveChat(ctx context.Context,
	in *proto.UserInGroup) (*proto.Result, error) {
	return h.db.UserUC.LeaveChat(in)
}

// AppendMessage to database
// to work correctly, specify the ID of the chat(in the message) in which
// the operation occurs
// Return id for this message, save it. It must be transferred to any message
// operations
func (h *Handler) AppendMessage(ctx context.Context,
	in *proto.Message) (*proto.MessageID, error) {
	return h.db.MessageUC.AppendOne(in)
}

// AppendMessages to database
func (h *Handler) AppendMessages(ctx context.Context,
	in *proto.Messages) (*proto.MessagesID, error) {
	return h.db.MessageUC.AppendMany(in)
}

// UpdateMessage in database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (h *Handler) UpdateMessage(ctx context.Context,
	in *proto.Message) (*proto.Result, error) {
	return h.db.MessageUC.Update(in)
}

// DeleteMessage from database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (h *Handler) DeleteMessage(ctx context.Context,
	in *proto.Message) (*proto.Result, error) {
	return h.db.MessageUC.Delete(in)
}

// GetChatMessages get all messages from the chad with specified id
func (h *Handler) GetChatMessages(ctx context.Context,
	in *proto.ChatID) (*proto.Messages, error) {
	return h.db.MessageUC.GetAll(in)
}

// 158 -> 108
