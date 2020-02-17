package database

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/repository/database"
)

type Message struct {
	db      infrastructure.Database
	message chat.MessageRepository
	logger  infrastructure.Logger
	trace   infrastructure.ErrorTrace
	timeout time.Duration
}

func NewMessage(
	dbI infrastructure.Database,
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
	photo infrastructure.PhotoService,
	timeout time.Duration,
) (*Message, error) {
	// check database interface given
	if dbI == nil {
		return nil, errors.New(infrastructure.ErrNoDatabase)
	}
	// overriding nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}
	// overriding nil value of ErrorTrace
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}
	// overriding nil value of PhotoService
	if photo == nil {
		photo = new(infrastructure.PhotoServiceNil)
	}
	messageRep, err := database.NewMessage(
		dbI,
		photo,
		logger,
		trace,
	)
	if err != nil {
		return nil, err
	}
	return &Message{
		db:      dbI,
		logger:  logger,
		trace:   trace,
		message: messageRep,
		timeout: timeout,
	}, nil
}

// AppendOne append message to database
// to work correctly, specify the ID of the chat(in the message) in which
// the operation occurs
// Return id for this message, save it. It must be transferred to any message
// operations
func (uc *Message) AppendOne(
	c context.Context,
	message *models.Message,
) (int32, error) {
	if message == nil {
		return 0, uc.trace.New(models.ErrNoMessage)
	}

	if message.ChatID <= 0 {
		return 0, uc.trace.New(InvalidMessageChatID)
	}

	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	return uc.message.CreateOne(ctx, message)
}

// AppendMany append messages to database
func (uc *Message) AppendMany(
	c context.Context,
	messages *models.Messages,
) ([]int32, error) {
	if messages == nil {
		return nil, uc.trace.New(models.ErrNoMessages)
	}

	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	return uc.message.CreateMany(ctx, messages)
}

// Update message in database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (uc *Message) Update(
	c context.Context,
	message *models.Message,
) error {
	if err := uc.checkExistingMessage(message); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	return uc.message.Update(ctx, message)
}

// Delete message from database
// to work correctly, specify the ID of the message in which
// the operation occurs
func (uc *Message) Delete(
	c context.Context,
	message *models.Message,
) error {
	if err := uc.checkExistingMessage(message); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	return uc.message.Delete(ctx, message)
}

// GetAll get all messages from the chad with specified id
func (uc *Message) GetAll(
	c context.Context,
	chatID int32,
) ([]*models.Message, error) {
	if err := uc.checkChatID(chatID); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c, uc.timeout)
	defer cancel()

	return uc.message.GetAll(ctx, chatID)
}

func (uc *Message) checkChatID(chatID int32) error {
	if chatID <= 0 {
		return uc.trace.New(InvalidChatID)
	}
	return nil
}

func (uc *Message) checkExistingMessage(
	message *models.Message,
) error {
	if message == nil {
		return uc.trace.New(models.ErrNoMessage)
	}
	if message.ChatID <= 0 {
		return uc.trace.New(InvalidMessageChatID)
	}
	if message.ID <= 0 {
		return uc.trace.New(InvalidMessageID)
	}
	return nil
}

// 193 -> 172
