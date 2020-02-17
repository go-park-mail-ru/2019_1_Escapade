package database

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/repository/database"
)

// User implements the interface UserUseCaseI
type User struct {
	db             infrastructure.Database
	trace          infrastructure.ErrorTrace
	userDB         api.UserRepositoryI
	recordDB       api.RecordRepositoryI
	contextTimeout time.Duration
}

// NewUser create new instance of User
func NewUser(
	dbI infrastructure.Database,
	trace infrastructure.ErrorTrace,
	timeout time.Duration,
) (*User, error) {
	if dbI == nil {
		return nil, errors.New(ErrNoDatabase)
	}
	recordRep, err := database.NewRecord(dbI, trace)
	if err != nil {
		return nil, err
	}
	userRep, err := database.NewUser(dbI, trace)
	if err != nil {
		return nil, err
	}
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}
	return &User{
		db:             dbI,
		trace:          trace,
		userDB:         userRep,
		recordDB:       recordRep,
		contextTimeout: timeout,
	}, nil
}

// CreateAccount check sql-injections and is name unique
// Then add cookie to database and returns session_id
func (usecase *User) CreateAccount(
	c context.Context,
	user *models.UserPrivateInfo,
) (int, error) {
	// check user not nil
	if user == nil {
		return 0, usecase.trace.New(InvalidUser)
	}
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()

	var (
		userID int
		err    error
		tx     infrastructure.Transaction
	)

	if tx, err = usecase.db.Begin(); err != nil {
		return userID, err
	}
	defer tx.Rollback()

	userTX, err := database.NewUser(tx, usecase.trace)
	if err != nil {
		return userID, err
	}
	recordTX, err := database.NewRecord(tx, usecase.trace)
	if err != nil {
		return userID, err
	}

	if userID, err = userTX.Create(ctx, user); err != nil {
		return userID, err
	}
	if err = recordTX.Create(ctx, userID); err != nil {
		return userID, err
	}

	err = tx.Commit()
	return userID, err
}

// EnterAccount check sql-injections and is password right
// Then add cookie to database and returns session_id
func (usecase *User) EnterAccount(
	c context.Context,
	name, password string,
) (int32, error) {
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()
	userID, _, err := usecase.userDB.CheckNamePassword(
		ctx,
		name,
		password,
	)
	return userID, err
}

// UpdateAccount gets name of Player from
// relation Session, cause we know that user has session
func (usecase *User) UpdateAccount(
	c context.Context,
	userID int32,
	user *models.UserPrivateInfo,
) error {
	// check user not nil
	if user == nil {
		return usecase.trace.New(InvalidUser)
	}
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()

	var (
		confirmedUser *models.UserPrivateInfo
		tx            infrastructure.Transaction
		err           error
	)

	if tx, err = usecase.db.Begin(); err != nil {
		return err
	}
	defer tx.Rollback()

	userTX, err := database.NewUser(tx, usecase.trace)
	if err != nil {
		return err
	}

	confirmedUser, err = userTX.FetchNamePassword(ctx, userID)
	if err != nil {
		return err
	}

	confirmedUser.Update(user)

	// TODO сделать разлогин других сессий юзера при смене пароля
	err = userTX.UpdateNamePassword(ctx, user)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// DeleteAccount deletes account
func (usecase *User) DeleteAccount(
	c context.Context,
	user *models.UserPrivateInfo,
) error {
	if user == nil {
		return usecase.trace.New(InvalidUser)
	}
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()

	var (
		err error
		tx  infrastructure.Transaction
	)

	if tx, err = usecase.db.Begin(); err != nil {
		return err
	}
	defer tx.Rollback()

	userTX, err := database.NewUser(tx, usecase.trace)
	if err != nil {
		return err
	}

	if err = userTX.Delete(ctx, user); err != nil {
		return err
	}

	// TODO delete all tokens
	// if err = db.deleteAllUserSessions(tx, user.Name); err != nil {
	// 	return
	// }

	err = tx.Commit()
	return err
}

// FetchAll get users
func (usecase *User) FetchAll(
	c context.Context,
	difficult, page, perPage int,
	sort string,
) ([]*models.UserPublicInfo, error) {
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()
	var (
		offset int
		limit  int
	)

	pageusers := 10 // в конфиг
	limit = perPage
	offset = limit * (page - 1)
	if offset > pageusers {
		return nil, usecase.trace.New("offset > pageusers")
	}
	if offset+limit >= pageusers {
		limit = pageusers - offset
		if limit == 0 {
			return nil, usecase.trace.New("pageusers - offset = 0")
		}
	}

	params := api.UsersSelectParams{
		Difficult: difficult,
		Offset:    offset,
		Limit:     limit,
		Sort:      sort,
	}

	players, err := usecase.userDB.FetchAll(ctx, params)
	if err != nil {
		return nil, err
	}

	return players, err
}

// FetchOne get one user
func (usecase *User) FetchOne(
	c context.Context,
	userID int32,
	difficult int,
) (*models.UserPublicInfo, error) {
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()

	user, err := usecase.userDB.FetchOne(
		ctx,
		userID,
		difficult,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// PagesCount returns amount of rows in table Player
// deleted on amount of rows in one page
func (usecase *User) PagesCount(
	c context.Context,
	perPage int,
) (int, error) {
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()
	count, err := usecase.userDB.PagesCount(ctx, perPage)
	return count, err
}
