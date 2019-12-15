package clients

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

//go:generate $GOPATH/bin/mockery -name "ChatI"

// ChatI interface of chat
type ChatI interface {
	Init(consul server.ConsulServiceI, required config.RequiredService) error
	Close() error

	CreateChat(in *proto.ChatWithUsers) (*proto.ChatID, error)
	GetChat(in *proto.Chat) (*proto.ChatID, error)

	AppendMessage(in *proto.Message) (*proto.MessageID, error)
	AppendMessages(in *proto.Messages) (*proto.MessagesID, error)
	UpdateMessage(in *proto.Message) (*proto.Result, error)
	DeleteMessage(in *proto.Message) (*proto.Result, error)
	GetChatMessages(in *proto.ChatID) (*proto.Messages, error)

	InviteToChat(in *proto.UserInGroup) (*proto.Result, error)
}
