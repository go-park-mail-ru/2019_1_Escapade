package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/return_errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
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
func (db *ChatUseCase) Create(chat *proto.ChatWithUsers) (*proto.ChatID, error) {
	if chat == nil {
		return &proto.ChatID{}, re.InvalidMessage()
	}

	var (
		tx     database.TransactionI
		err    error
		ChatID = &proto.ChatID{}
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
func (db *ChatUseCase) GetOne(chat *proto.Chat) (*proto.ChatID, error) {
	if chat == nil {
		return &proto.ChatID{}, re.InvalidChatID()
	}
	var (
		tx  database.TransactionI
		err error
	)
	if tx, err = db.Db.Begin(); err != nil {
		return nil, err
	}
	defer tx.Rollback()

	id, err := db.chat.get(tx, chat)
	if err != nil {

		if id, err = db.chat.create(tx, chat.Type, chat.TypeId); err != nil {
			return id, err
		}
	}
	if err = tx.Commit(); err != nil {
		return id, err
	}
	return id, err
}
