package chat

import (
	context "context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
)

type ChatRepositoryI interface {
	create(tx database.TransactionI, chatType ChatType, typeID int32) (*ChatID, error)
	get(Db database.DatabaseI, chat *Chat) (*ChatID, error)
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

	Create(ctx context.Context, chat *ChatWithUsers) (*ChatID, error)
	GetOne(ctx context.Context, chat *Chat) (*ChatID, error)
}

type UserUseCaseI interface {
	database.UserCaseI

	InviteToChat(ctx context.Context, userInChat *UserInGroup) (*Result, error)
	LeaveChat(ctx context.Context, userInChat *UserInGroup) (*Result, error)
}

type MessageUseCaseI interface {
	database.UserCaseI

	AppendOne(ctx context.Context, message *Message) (*MessageID, error)
	AppendMany(ctx context.Context, messages *Messages) (*MessagesID, error)

	Update(ctx context.Context, message *Message) (*Result, error)
	Delete(ctx context.Context, message *Message) (*Result, error)

	GetAll(ctx context.Context, chatID *ChatID) (*Messages, error)
}
