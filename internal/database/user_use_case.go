package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// UserUseCase implements the interface UserUseCaseI
type UserUseCase struct {
	UseCaseBase
	user   UserRepositoryI
	record RecordRepositoryI
}

func (db *UserUseCase) Init(user UserRepositoryI, record RecordRepositoryI) {
	db.user = user
	db.record = record
}

// CreateAccount check sql-injections and is name unique
// Then add cookie to database and returns session_id
func (db *UserUseCase) CreateAccount(user *models.UserPrivateInfo) (int, error) {

	var (
		userID int
		err    error
		tx     transactionI
	)

	if tx, err = db.Db.Begin(); err != nil {
		return userID, err
	}
	defer tx.Rollback()

	if userID, err = db.user.create(tx, user); err != nil {
		return userID, err
	}

	if err = db.record.create(tx, userID); err != nil {
		return userID, err
	}

	err = tx.Commit()
	return userID, err
}

// EnterAccount check sql-injections and is password right
// Then add cookie to database and returns session_id
func (db *UserUseCase) EnterAccount(name, password string) (int32, error) {

	var (
		tx     transactionI
		userID int32
		err    error
	)

	if tx, err = db.Db.Begin(); err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if userID, _, err = db.user.checkNamePassword(tx, name, password); err != nil {
		return userID, err
	}

	err = tx.Commit()
	return userID, err
}

// UpdateAccount gets name of Player from
// relation Session, cause we know that user has session
func (db *UserUseCase) UpdateAccount(userID int32, user *models.UserPrivateInfo) error {
	var (
		confirmedUser *models.UserPrivateInfo
		tx            transactionI
		err           error
	)

	if tx, err = db.Db.Begin(); err != nil {
		return err
	}
	defer tx.Rollback()

	if confirmedUser, err = db.user.fetchNamePassword(tx, userID); err != nil {
		return err
	}

	confirmedUser.Update(user)

	if err = db.user.updateNamePassword(tx, user); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// DeleteAccount deletes account
func (db *UserUseCase) DeleteAccount(user *models.UserPrivateInfo) (err error) {
	var tx transactionI

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if err = db.user.delete(tx, user); err != nil {
		return
	}

	// TODO delete all tokens
	// if err = db.deleteAllUserSessions(tx, user.Name); err != nil {
	// 	return
	// }

	err = tx.Commit()
	return
}

// FetchAll get users
func (db *UserUseCase) FetchAll(difficult int, page int, perPage int,
	sort string) (players []*models.UserPublicInfo, err error) {
	var (
		offset int
		limit  int
		tx     transactionI
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	pageusers := 10 // в конфиг
	limit = perPage
	offset = limit * (page - 1)
	if offset > pageusers {
		return
	}
	if offset+limit >= pageusers {
		limit = pageusers - offset
		if limit == 0 {
			return
		}
	}

	params := UsersSelectParams{
		Difficult: difficult,
		Offset:    offset,
		Limit:     limit,
		Sort:      sort,
	}

	if players, err = db.user.fetchAll(tx, params); err != nil {
		return
	}

	err = tx.Commit()
	return
}

// FetchOne get one user
func (db *UserUseCase) FetchOne(userID int32, difficult int) (user *models.UserPublicInfo, err error) {
	var tx transactionI

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if user, err = db.user.fetchOne(tx, userID, difficult); err != nil {
		return
	}

	err = tx.Commit()
	return
}

// PagesCount returns amount of rows in table Player
// deleted on amount of rows in one page
func (db *UserUseCase) PagesCount(perPage int) (int, error) {
	return db.user.pagesCount(db.Db, perPage)
}
