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

type User struct {
	db             infrastructure.Database
	user           chat.UserRepository
	logger         infrastructure.Logger
	trace          infrastructure.ErrorTrace
	contextTimeout time.Duration
}

func NewUser(
	dbI infrastructure.Database,
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
	timeout time.Duration,
) (*User, error) {
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
	userRep, err := database.NewUser(dbI, logger, trace)
	if err != nil {
		return nil, err
	}
	return &User{
		db:             dbI,
		logger:         logger,
		trace:          trace,
		user:           userRep,
		contextTimeout: timeout,
	}, nil
}

// InviteToChat invite user to the chat
// to work correctly, specify user and id of the chat
func (uc *User) InviteToChat(
	c context.Context,
	userInChat *models.UserInGroup,
) error {
	if err := uc.check(userInChat); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(
		c,
		uc.contextTimeout,
	)
	defer cancel()

	err := uc.user.Create(
		ctx,
		userInChat.Chat.ID,
		userInChat.User,
	)

	return err
}

// LeaveChat leave user from the chat
// to work correctly, specify user and id of the chat
func (uc *User) LeaveChat(
	c context.Context,
	userInChat *models.UserInGroup,
) error {

	if err := uc.check(userInChat); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(
		c,
		uc.contextTimeout,
	)
	defer cancel()

	return uc.user.Delete(ctx, userInChat)
}

func (uc *User) check(userInChat *models.UserInGroup) error {
	if userInChat == nil {
		return uc.trace.New(models.ErrNoUserInGroup)
	}

	if userInChat.User.ID <= 0 {
		return uc.trace.New(InvalidUserID)
	}

	if userInChat.Chat.ID <= 0 {
		return uc.trace.New(InvalidUserChatID)
	}
	return nil
}
