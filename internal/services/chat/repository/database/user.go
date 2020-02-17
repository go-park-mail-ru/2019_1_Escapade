package database

import (
	"context"
	"errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/models"
)

type User struct {
	db     infrastructure.Execer
	logger infrastructure.Logger
	trace  infrastructure.ErrorTrace
}

func NewUser(
	db infrastructure.Execer,
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
) (*User, error) {
	// check database interface given
	if db == nil {
		return nil, errors.New(ErrNoDatabase)
	}
	// overriding nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}
	// overriding nil value of ErrorTrace
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}
	return &User{
		db:     db,
		logger: logger,
		trace:  trace,
	}, nil
}

func addUserToQuery(user *models.User) string {
	return "('" + utils.String32(user.ID) + "',$1)"
}

func (rep *User) Create(
	ctx context.Context,
	chatID int32,
	users ...*models.User,
) error {
	var (
		err error
	)
	sqlInsert := `INSERT INTO UserInChat(user_id, chat_id) VALUES `

	if len(users) == 0 {
		return nil
	}
	for i, user := range users {
		if i == 0 {
			sqlInsert += addUserToQuery(user)
		} else {
			sqlInsert += "," + addUserToQuery(user)
		}
	}

	_, err = rep.db.ExecContext(ctx, sqlInsert, chatID)

	return err
}

func (rep *User) Delete(
	ctx context.Context,
	userInGroup *models.UserInGroup,
) error {
	var (
		id        int32
		err       error
		sqlDelete = `	
		Delete from UserInChat 
			where user_id = $1 and chat_id = $2;
		`
		row = rep.db.QueryRowContext(
			ctx,
			sqlDelete,
			userInGroup.User.ID,
			userInGroup.Chat.ID,
		)
	)

	if err = row.Scan(&id); err != nil {
		rep.logger.Println("cant delete message", err)
		return err
	}

	return nil
}
