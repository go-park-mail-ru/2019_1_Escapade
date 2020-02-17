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

type Chat struct {
	db             infrastructure.Database
	logger         infrastructure.Logger
	trace          infrastructure.ErrorTrace
	chat           chat.ChatRepository
	user           chat.UserRepository
	contextTimeout time.Duration
}

func NewChat(
	dbI infrastructure.Database,
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
	timeout time.Duration,
) (*Chat, error) {
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
	chatRep, err := database.NewChat(dbI, logger, trace)
	if err != nil {
		return nil, err
	}
	userRep, err := database.NewUser(dbI, logger, trace)
	if err != nil {
		return nil, err
	}
	return &Chat{
		db:             dbI,
		logger:         logger,
		trace:          trace,
		chat:           chatRep,
		user:           userRep,
		contextTimeout: timeout,
	}, nil
}

// Create chat with or without users.
// Specify the type of chat and id received from the corresponding database table
// Return id for this chat, save it. It must be transferred to any chat operations
func (uc *Chat) Create(
	c context.Context,
	chat *models.ChatWithUsers,
) (int32, error) {
	if chat == nil {
		return 0, uc.trace.New(models.ErrNoChatWithUsers)
	}

	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	tx, err := uc.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	chatTX, err := database.NewChat(tx, uc.logger, uc.trace)
	if err != nil {
		return 0, err
	}

	ChatID, err := chatTX.Create(ctx, chat.Type, chat.TypeID)
	if err != nil {
		return 0, err
	}

	userTX, err := database.NewUser(tx, uc.logger, uc.trace)
	if err != nil {
		return 0, err
	}

	err = userTX.Create(ctx, ChatID, chat.Users...)
	if err != nil {
		return ChatID, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return ChatID, nil
}

// GetOne get the ID of the chat, based on its type and the passed ID of this type
func (uc *Chat) GetOne(
	c context.Context,
	chat *models.Chat,
) (int32, error) {
	if chat == nil {
		return 0, uc.trace.New(models.ErrNoChat)
	}
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	id, err := uc.chat.Get(ctx, chat)
	if err != nil {
		id, err = uc.chat.Create(ctx, chat.Type, chat.TypeID)
	}
	return id, err
}
