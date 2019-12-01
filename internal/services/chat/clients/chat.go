package clients

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/clients"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

type Chat struct {
	clients.BaseService
	chat     proto.ChatServiceClient
	premove  func()
	postmove func(error)
}

func (c *Chat) Init(consul *server.ConsulService, required config.RequiredService) error {
	err := c.BaseService.Init(consul, required)
	if err != nil {
		return err
	}
	c.chat = proto.NewChatServiceClient(c.GrcpConn)
	c.premove = func() {}
	c.postmove = func(err error) {
		if err != nil {
			c.ErrorIncrese()
		}
	}
	return nil
}

func (c *Chat) do(f func(context.Context) error) {
	c.premove()
	err := f(context.Background())
	c.postmove(err)
}

// CreateChat chat with or without users.
// Specify the type of chat and id received from the corresponding
//  database table
// Return id for this chat, save it. It must be transferred to any
//  chat operations
func (c *Chat) CreateChat(in *proto.ChatWithUsers) (*proto.ChatID, error) {
	var (
		out *proto.ChatID
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.CreateChat(ctx, in)
		return err
	})
	return out, err
}

// GetChat get the ID of the chat, based on its type and the passed ID of this type
func (c *Chat) GetChat(in *proto.Chat) (*proto.ChatID, error) {
	var (
		out *proto.ChatID
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.GetChat(ctx, in)
		return err
	})
	return out, err
}

// InviteToChat invite user to the chat
// to work correctly, specify user and id of the chat
func (c *Chat) InviteToChat(in *proto.UserInGroup) (*proto.Result, error) {
	var (
		out *proto.Result
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.InviteToChat(ctx, in)
		return err
	})
	return out, err
}

// LeaveChat leave user from the chat
// to work correctly, specify user and id of the chat
func (c *Chat) LeaveChat(in *proto.UserInGroup) (*proto.Result, error) {
	var (
		out *proto.Result
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.LeaveChat(ctx, in)
		return err
	})
	return out, err

}

// AppendMessage to database
// to work correctly, specify the ID of the chat(in the message) in which
// the operation occurs
// Return id for this message, save it. It must be transferred to any message
// operations
func (c *Chat) AppendMessage(in *proto.Message) (*proto.MessageID, error) {
	var (
		out *proto.MessageID
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.AppendMessage(ctx, in)
		return err
	})
	return out, err
}

// AppendMessages to database
func (c *Chat) AppendMessages(in *proto.Messages) (*proto.MessagesID, error) {
	var (
		out *proto.MessagesID
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.AppendMessages(ctx, in)
		return err
	})
	return out, err
}

// UpdateMessage in database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (c *Chat) UpdateMessage(in *proto.Message) (*proto.Result, error) {
	var (
		out *proto.Result
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.UpdateMessage(ctx, in)
		return err
	})
	return out, err
}

// DeleteMessage from database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (c *Chat) DeleteMessage(in *proto.Message) (*proto.Result, error) {
	var (
		out *proto.Result
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.DeleteMessage(ctx, in)
		return err
	})
	return out, err
}

// GetChatMessages get all messages from the chad with specified id
func (c *Chat) GetChatMessages(in *proto.ChatID) (*proto.Messages, error) {
	var (
		out *proto.Messages
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.GetChatMessages(ctx, in)
		return err
	})
	return out, err

}
