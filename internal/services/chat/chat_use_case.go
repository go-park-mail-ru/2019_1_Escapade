package chat

import (
	context "context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

type ChatUseCase struct {
	database.UseCaseBase
	chat ChatRepositoryI
	user UserRepositoryI
}

func (db *ChatUseCase) Init(chat ChatRepositoryI, user UserRepositoryI) {
	db.chat = chat
	db.user = user
}

// Create chat with or without users.
// Specify the type of chat and id received from the corresponding database table
// Return id for this chat, save it. It must be transferred to any chat operations
func (db *ChatUseCase) Create(ctx context.Context, chat *ChatWithUsers) (*ChatID, error) {
	if chat == nil {
		return &ChatID{}, re.InvalidMessage()
	}

	var (
		tx     database.TransactionI
		err    error
		ChatID = &ChatID{}
	)

	if tx, err = db.Db.Begin(); err != nil {
		return ChatID, err
	}
	defer tx.Rollback()

	if ChatID, err = db.chat.create(tx, chat.Type, chat.TypeId); err != nil {
		return ChatID, err
	}

	if err = db.user.create(tx, ChatID.Value, chat.Users...); err != nil {
		return ChatID, err
	}

	if err = tx.Commit(); err != nil {
		return ChatID, err
	}
	return ChatID, err
}

// GetOne get the ID of the chat, based on its type and the passed ID of this type
func (db *ChatUseCase) GetOne(ctx context.Context, chat *Chat) (*ChatID, error) {
	if chat == nil {
		return &ChatID{}, re.InvalidChatID()
	}

	id, err := db.chat.get(db.Db, chat)
	if err != nil {
		var (
			tx database.TransactionI
		)
		if tx, err = db.Db.Begin(); err != nil {
			return id, err
		}
		defer tx.Rollback()
		if id, err = db.chat.create(tx, chat.Type, chat.TypeId); err != nil {
			return id, err
		}
		if err = tx.Commit(); err != nil {
			return id, err
		}
	}

	return id, err
}
