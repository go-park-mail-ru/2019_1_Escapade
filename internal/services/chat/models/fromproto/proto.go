package fromproto

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
	"github.com/golang/protobuf/ptypes"
)

type User proto.User

func (user *User) Proto() *proto.User {
	return (*proto.User)(user)
}

func (user *User) Get() *models.User {
	if user == nil {
		return nil
	}
	return &models.User{
		ID:     user.Id,
		Name:   user.Name,
		Photo:  user.Photo,
		Status: int32(user.Status),
	}
}

func NewUser(m *models.User) *User {
	if m == nil {
		return nil
	}
	return &User{
		Id:     m.ID,
		Name:   m.Name,
		Photo:  m.Photo,
		Status: proto.Status(m.Status),
	}
}

type Message proto.Message

func (message *Message) Proto() *proto.Message {
	return (*proto.Message)(message)
}

func (message *Message) Get() *models.Message {
	if message == nil {
		return nil
	}
	pTime, err := ptypes.Timestamp(message.Time)
	if err != nil {
		return nil
	}
	return &models.Message{
		ID:     message.Id,
		Answer: (*Message)(message.Answer).Get(),
		Text:   message.Text,
		From:   (*User)(message.From).Get(),
		To:     (*User)(message.To).Get(),
		ChatID: message.ChatId,
		Time:   pTime,
		Edited: message.Edited,
	}
}

func NewMessage(m *models.Message) *Message {
	if m == nil {
		return nil
	}
	pTime, err := ptypes.TimestampProto(m.Time)
	if err != nil {
		return nil
	}
	return &Message{
		Id:     m.ID,
		Answer: NewMessage(m.Answer).Proto(),
		Text:   m.Text,
		From:   NewUser(m.From).Proto(),
		To:     NewUser(m.To).Proto(),
		ChatId: m.ChatID,
		Time:   pTime,
		Edited: m.Edited,
	}
}

type Messages proto.Messages

func (messages *Messages) Get() *models.Messages {
	if messages == nil {
		return nil
	}

	var msgs = make([]*models.Message, len(messages.Messages))
	for i, msg := range messages.Messages {
		msgs[i] = (*Message)(msg).Get()
	}
	return &models.Messages{
		Messages:    msgs,
		BlockSize:   messages.BlockSize,
		BlockAmount: messages.BlockAmount,
		BlockNumber: messages.BlockNumber,
	}
}

func (messages *Messages) Proto() *proto.Messages {
	return (*proto.Messages)(messages)
}

func NewMessages(m *models.Messages) *Messages {
	if m == nil {
		return nil
	}
	var msgs = make([]*proto.Message, len(m.Messages))
	for i, msg := range m.Messages {
		msgs[i] = NewMessage(msg).Proto()
	}
	return &Messages{
		Messages:    msgs,
		BlockSize:   m.BlockSize,
		BlockAmount: m.BlockAmount,
		BlockNumber: m.BlockNumber,
	}
}

type Chat proto.Chat

func (chat *Chat) Get() *models.Chat {
	if chat == nil {
		return nil
	}

	var msgs = make([]*models.Messages, len(chat.Messages))
	for i, msg := range chat.Messages {
		msgs[i] = (*Messages)(msg).Get()
	}
	return &models.Chat{
		Messages: msgs,
		ID:       chat.Id,
		Type:     chat.Type,
		TypeID:   chat.TypeId,
	}
}

func (chat *Chat) Proto() *proto.Chat {
	return (*proto.Chat)(chat)
}

func NewChat(m *models.Chat) *Chat {
	if m == nil {
		return nil
	}
	var msgs = make([]*proto.Messages, len(m.Messages))
	for i, msg := range m.Messages {
		msgs[i] = NewMessages(msg).Proto()
	}
	return &Chat{
		Messages: msgs,
		Id:       m.ID,
		Type:     m.Type,
		TypeId:   m.TypeID,
	}
}

type UserInGroup proto.UserInGroup

func (userInGroup *UserInGroup) Get() *models.UserInGroup {
	if userInGroup == nil {
		return nil
	}
	return &models.UserInGroup{
		User: (*User)(userInGroup.User).Get(),
		Chat: (*Chat)(userInGroup.Chat).Get(),
	}
}

func (userInGroup *UserInGroup) Proto() *proto.UserInGroup {
	return (*proto.UserInGroup)(userInGroup)
}

func NewUserInGroup(m *models.UserInGroup) *UserInGroup {
	if m == nil {
		return nil
	}
	return &UserInGroup{
		User: NewUser(m.User).Proto(),
		Chat: NewChat(m.Chat).Proto(),
	}
}

type ChatWithUsers proto.ChatWithUsers

func (chat *ChatWithUsers) Get() *models.ChatWithUsers {
	if chat == nil {
		return nil
	}
	var users = make([]*models.User, len(chat.Users))
	for i, user := range chat.Users {
		users[i] = (*User)(user).Get()
	}
	return &models.ChatWithUsers{
		Type:   chat.Type,
		TypeID: chat.TypeId,
		Users:  users,
	}
}

func (chat *ChatWithUsers) Proto() *proto.ChatWithUsers {
	return (*proto.ChatWithUsers)(chat)
}

func NewChatWithUsers(m *models.ChatWithUsers) *ChatWithUsers {
	if m == nil {
		return nil
	}
	var users = make([]*proto.User, len(m.Users))
	for i, user := range m.Users {
		users[i] = NewUser(user).Proto()
	}
	return &ChatWithUsers{
		Type:   m.Type,
		TypeId: m.TypeID,
		Users:  users,
	}
}
