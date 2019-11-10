package clients

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	s_chat "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat"
)

type Chat struct {
	BaseService
	chat     s_chat.ChatServiceClient
	premove  func()
	postmove func(error)
}

func (c *Chat) Init(consul *server.ConsulService, required config.RequiredService) error {
	err := c.BaseService.Init(consul, required)
	if err != nil {
		return err
	}
	c.chat = s_chat.NewChatServiceClient(c.grcpConn)
	c.premove = func() {}
	c.postmove = func(err error) {
		if err != nil {
			c.errorIncrese()
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
// Specify the type of chat and id received from the corresponding database table
// Return id for this chat, save it. It must be transferred to any chat operations
func (c *Chat) CreateChat(in *s_chat.ChatWithUsers) (*s_chat.ChatID, error) {
	var (
		out *s_chat.ChatID
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.CreateChat(ctx, in)
		return err
	})
	return out, err
}

// GetChat get the ID of the chat, based on its type and the passed ID of this type
func (c *Chat) GetChat(in *s_chat.Chat) (*s_chat.ChatID, error) {
	var (
		out *s_chat.ChatID
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
func (c *Chat) InviteToChat(in *s_chat.UserInGroup) (*s_chat.Result, error) {
	var (
		out *s_chat.Result
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
func (c *Chat) LeaveChat(in *s_chat.UserInGroup) (*s_chat.Result, error) {
	var (
		out *s_chat.Result
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
func (c *Chat) AppendMessage(in *s_chat.Message) (*s_chat.MessageID, error) {
	var (
		out *s_chat.MessageID
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.AppendMessage(ctx, in)
		return err
	})
	return out, err
}

// AppendMessages to database
func (c *Chat) AppendMessages(in *s_chat.Messages) (*s_chat.MessagesID, error) {
	var (
		out *s_chat.MessagesID
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
func (c *Chat) UpdateMessage(in *s_chat.Message) (*s_chat.Result, error) {
	var (
		out *s_chat.Result
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
func (c *Chat) DeleteMessage(in *s_chat.Message) (*s_chat.Result, error) {
	var (
		out *s_chat.Result
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.DeleteMessage(ctx, in)
		return err
	})
	return out, err
}

// GetChatMessages get all messages from the chad with specified id
func (c *Chat) GetChatMessages(in *s_chat.ChatID) (*s_chat.Messages, error) {
	var (
		out *s_chat.Messages
		err error
	)
	c.do(func(ctx context.Context) error {
		out, err = c.chat.GetChatMessages(ctx, in)
		return err
	})
	return out, err

}
