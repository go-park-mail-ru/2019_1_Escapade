package chat

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
)

type MessageUseCase struct {
	database.UseCaseBase
	message MessageRepositoryI
}

func (db *MessageUseCase) Init(message MessageRepositoryI) {
	db.message = message
}

// AppendOne append message to database
// to work correctly, specify the ID of the chat(in the message) in which
// the operation occurs
// Return id for this message, save it. It must be transferred to any message
// operations
func (db *MessageUseCase) AppendOne(message *Message) (*MessageID, error) {
	if message == nil {
		return &MessageID{}, re.InvalidMessage()
	}

	if message.ChatId <= 0 {
		return &MessageID{}, re.InvalidMessageID()
	}

	id, err := db.message.createOne(db.Db, message)
	return id, err
}

// AppendMany append messages to database
func (db *MessageUseCase) AppendMany(messages *Messages) (*MessagesID, error) {
	if messages == nil {
		return &MessagesID{}, re.InvalidMessage()
	}

	ids, err := db.message.createMany(db.Db, messages)
	return ids, err
}

// Update message in database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (db *MessageUseCase) Update(message *Message) (*Result, error) {
	if message == nil {
		return &Result{}, re.InvalidMessage()
	}
	if message.Id <= 0 {
		return &Result{}, re.InvalidMessageID()
	}

	res, err := db.message.update(db.Db, message)
	return res, err
}

// Delete message from database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (db *MessageUseCase) Delete(message *Message) (*Result, error) {
	if message == nil {
		return &Result{}, re.InvalidMessage()
	}
	if message.Id <= 0 {
		return &Result{}, re.InvalidMessageID()
	}

	res, err := db.message.delete(db.Db, message)
	return res, err
}

// GetChatMessages get all messages from the chad with specified id
func (db *MessageUseCase) GetAll(chatID *ChatID) (*Messages, error) {

	if chatID == nil {
		return &Messages{}, re.InvalidMessage()
	}

	messages, err := db.message.getAll(db.Db, chatID)
	return messages, err
}
