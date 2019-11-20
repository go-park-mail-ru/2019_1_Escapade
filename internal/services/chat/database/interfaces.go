package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

const (
	LobbyType = 0
	RoomType  = 1
	UserType  = 2
)

type ChatRepositoryI interface {
	create(tx database.TransactionI, chatType, typeID int32) (*proto.ChatID, error)
	get(tx database.TransactionI, chat *proto.Chat) (*proto.ChatID, error)
}

type UserRepositoryI interface {
	create(tx database.TransactionI, chatID int32, users ...*proto.User) error
	delete(Db database.DatabaseI, userInGroup *proto.UserInGroup) (*proto.Result, error)
}

type MessageRepositoryI interface {
	createOne(Db database.DatabaseI, message *proto.Message) (*proto.MessageID, error)
	createMany(Db database.DatabaseI, messages *proto.Messages) (*proto.MessagesID, error)

	update(Db database.DatabaseI, message *proto.Message) (*proto.Result, error)
	delete(Db database.DatabaseI, message *proto.Message) (*proto.Result, error)

	getAll(Db database.DatabaseI, chatID *proto.ChatID) (*proto.Messages, error)
}

type ChatUseCaseI interface {
	database.UserCaseI

	Create(chat *proto.ChatWithUsers) (*proto.ChatID, error)
	GetOne(chat *proto.Chat) (*proto.ChatID, error)
}

type UserUseCaseI interface {
	database.UserCaseI

	InviteToChat(userInChat *proto.UserInGroup) (*proto.Result, error)
	LeaveChat(userInChat *proto.UserInGroup) (*proto.Result, error)
}

type MessageUseCaseI interface {
	database.UserCaseI

	AppendOne(message *proto.Message) (*proto.MessageID, error)
	AppendMany(messages *proto.Messages) (*proto.MessagesID, error)

	Update(message *proto.Message) (*proto.Result, error)
	Delete(message *proto.Message) (*proto.Result, error)

	GetAll(chatID *proto.ChatID) (*proto.Messages, error)
}
