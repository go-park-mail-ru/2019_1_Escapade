package grpc

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/models/fromproto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

type ChatServiceServer struct {
	chat    chat.ChatUseCase
	user    chat.UserUseCase
	message chat.MessageUseCase
}

// Init initialize Handler
func NewChatServiceServer(
	chatUseCase chat.ChatUseCase,
	userUseCase chat.UserUseCase,
	messageUseCase chat.MessageUseCase,
) (*ChatServiceServer, error) {
	if chatUseCase == nil {
		return nil, errors.New(chat.ErrNoChatUseCase)
	}
	if userUseCase == nil {
		return nil, errors.New(chat.ErrNoUserUseCase)
	}
	if messageUseCase == nil {
		return nil, errors.New(chat.ErrNoMessageUseCase)
	}
	return &ChatServiceServer{
		chat: chatUseCase,
	}, nil
}

// CreateChat chat with or without users.
// Specify the type of chat and id received from the corresponding database
// table. Return id for this chat, save it. It must be transferred to any
// chat operations
func (server *ChatServiceServer) CreateChat(
	ctx context.Context,
	in *proto.ChatWithUsers,
) (*proto.ChatID, error) {
	id, err := server.chat.Create(
		ctx,
		(*fromproto.ChatWithUsers)(in).Get(),
	)
	return &proto.ChatID{Value: id}, err
}

// GetChat get the ID of the chat, based on its type and the passed ID of
// this type
func (server *ChatServiceServer) GetChat(
	ctx context.Context,
	in *proto.Chat,
) (*proto.ChatID, error) {
	id, err := server.chat.GetOne(
		ctx,
		(*fromproto.Chat)(in).Get(),
	)
	return &proto.ChatID{Value: id}, err
}

// InviteToChat invite user to the chat
// to work correctly, specify user and id of the chat
func (server *ChatServiceServer) InviteToChat(
	ctx context.Context,
	in *proto.UserInGroup,
) (*proto.Result, error) {
	err := server.user.InviteToChat(
		ctx,
		(*fromproto.UserInGroup)(in).Get(),
	)
	return &proto.Result{Done: err == nil}, err
}

// LeaveChat leave user from the chat
// to work correctly, specify user and id of the chat
func (server *ChatServiceServer) LeaveChat(
	ctx context.Context,
	in *proto.UserInGroup,
) (*proto.Result, error) {
	err := server.user.LeaveChat(
		ctx,
		(*fromproto.UserInGroup)(in).Get(),
	)
	return &proto.Result{Done: err == nil}, err
}

// AppendMessage to database
// to work correctly, specify the ID of the chat(in the message) in which
// the operation occurs
// Return id for this message, save it. It must be transferred to any message
// operations
func (server *ChatServiceServer) AppendMessage(
	ctx context.Context,
	in *proto.Message,
) (*proto.MessageID, error) {
	id, err := server.message.AppendOne(
		ctx,
		(*fromproto.Message)(in).Get(),
	)
	return &proto.MessageID{Value: id}, err
}

// AppendMessages to database
func (server *ChatServiceServer) AppendMessages(
	ctx context.Context,
	in *proto.Messages,
) (*proto.MessagesID, error) {
	ids, err := server.message.AppendMany(
		ctx,
		(*fromproto.Messages)(in).Get(),
	)
	protoIDs := make([]*proto.MessageID, len(ids))
	for i, id := range ids {
		protoIDs[i] = &proto.MessageID{Value: id}
	}
	return &proto.MessagesID{Values: protoIDs}, err
}

// UpdateMessage in database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (server *ChatServiceServer) UpdateMessage(
	ctx context.Context,
	in *proto.Message,
) (*proto.Result, error) {
	err := server.message.Update(
		ctx,
		(*fromproto.Message)(in).Get(),
	)
	return &proto.Result{Done: err == nil}, err
}

// DeleteMessage from database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (server *ChatServiceServer) DeleteMessage(
	ctx context.Context,
	in *proto.Message,
) (*proto.Result, error) {
	err := server.message.Delete(
		ctx,
		(*fromproto.Message)(in).Get(),
	)
	return &proto.Result{Done: err == nil}, err
}

// GetChatMessages get all messages from the chad with specified id
func (server *ChatServiceServer) GetChatMessages(
	ctx context.Context,
	in *proto.ChatID,
) (*proto.Messages, error) {
	msgs, err := server.message.GetAll(
		ctx,
		in.Value,
	)
	protoMessages := make([]*proto.Message, len(msgs))
	for i, msg := range msgs {
		protoMessages[i] = fromproto.NewMessage(msg).Proto()
	}
	//TODO implement BlockSize,BlockAmount,BlockNumber
	return &proto.Messages{
		Messages: protoMessages,
	}, err
}

// 158 -> 108 -> 171
