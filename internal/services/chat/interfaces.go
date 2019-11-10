package chat

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
)

const (
	LobbyType = 0
	RoomType  = 1
	UserType  = 2
)

type ChatRepositoryI interface {
	create(tx database.TransactionI, chatType, typeID int32) (*ChatID, error)
	get(tx database.TransactionI, chat *Chat) (*ChatID, error)
}

type UserRepositoryI interface {
	create(tx database.TransactionI, chatID int32, users ...*User) error
	delete(Db database.DatabaseI, userInGroup *UserInGroup) (*Result, error)
}

type MessageRepositoryI interface {
	createOne(Db database.DatabaseI, message *Message) (*MessageID, error)
	createMany(Db database.DatabaseI, messages *Messages) (*MessagesID, error)

	update(Db database.DatabaseI, message *Message) (*Result, error)
	delete(Db database.DatabaseI, message *Message) (*Result, error)

	getAll(Db database.DatabaseI, chatID *ChatID) (*Messages, error)
}

type ChatUseCaseI interface {
	database.UserCaseI

	Create(chat *ChatWithUsers) (*ChatID, error)
	GetOne(chat *Chat) (*ChatID, error)
}

type UserUseCaseI interface {
	database.UserCaseI

	InviteToChat(userInChat *UserInGroup) (*Result, error)
	LeaveChat(userInChat *UserInGroup) (*Result, error)
}

type MessageUseCaseI interface {
	database.UserCaseI

	AppendOne(message *Message) (*MessageID, error)
	AppendMany(messages *Messages) (*MessagesID, error)

	Update(message *Message) (*Result, error)
	Delete(message *Message) (*Result, error)

	GetAll(chatID *ChatID) (*Messages, error)
}
