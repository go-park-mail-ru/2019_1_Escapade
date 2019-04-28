package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"database/sql"
	"fmt"
)

// Register check sql-injections and are email and name unique
// Then add cookie to database and returns session_id
func (db *DataBase) Register(user *models.UserPrivateInfo, sessionID string) (userID int, err error) {

	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if userID, err = db.createPlayer(tx, user); err != nil {
		err = re.ErrorUserIsExist()
		fmt.Println("database/register - fail creating User")
		return
	}

	if err = db.createSession(tx, userID, sessionID); err != nil {
		fmt.Println("database/register - fail creating Session")
		return
	}

	if err = db.createRecords(tx, userID); err != nil {
		fmt.Println("database/register - fail creating Session")
		return
	}

	fmt.Println("database/register +")

	err = tx.Commit()
	return
}

// Login check sql-injections and is password right
// Then add cookie to database and returns session_id
func (db *DataBase) Login(user *models.UserPrivateInfo, sessionID string) (found *models.UserPublicInfo, err error) {

	var (
		tx     *sql.Tx
		userID int
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if userID, found, err = db.checkBunch(tx, user.Email, user.Password); err != nil {
		fmt.Println("database/login - fail enter")
		return
	}

	if err = db.createSession(tx, userID, sessionID); err != nil {
		fmt.Println("database/login - fail creating Session")
		return
	}

	if err = db.updatePlayerLastSeen(tx, userID); err != nil {
		fmt.Println("database/login - fail updatePlayerLastSeen")
		return
	}

	fmt.Println("database/login +")

	err = tx.Commit()
	return
}

// UpdatePlayerPersonalInfo gets name of Player from
// relation Session, cause we know that user has session
func (db *DataBase) UpdatePlayerPersonalInfo(userID int, user *models.UserPrivateInfo) (err error) {
	var (
		confirmedUser *models.UserPrivateInfo
		tx            *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if confirmedUser, err = db.getPrivateInfo(tx, userID); err != nil {
		return
	}

	confirmedUser.Update(user)

	if err = db.updatePlayerPersonalInfo(tx, user); err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *DataBase) GetUsers(difficult int, page int, perPage int,
	sort string) (players []*models.UserPublicInfo, err error) {
	var (
		offset int
		limit  int
		tx     *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	limit = perPage
	offset = limit * (page - 1)
	if offset > db.PageUsers {
		return
	}
	if offset+limit >= db.PageUsers {
		limit = db.PageUsers - offset
		if limit == 0 {
			return
		}
	}

	if players, err = db.getUsers(tx, difficult, offset, limit, sort); err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *DataBase) GetUser(userID int, difficult int) (user *models.UserPublicInfo, err error) {

	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if user, err = db.getUser(tx, userID, difficult); err != nil {
		return
	}

	err = tx.Commit()
	return
}

// DeleteAccount deletes account
func (db *DataBase) DeleteAccount(user *models.UserPrivateInfo) (err error) {

	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if err = db.deletePlayer(tx, user); err != nil {
		fmt.Println("database/DeleteAccount - fail deletting User")
		return
	}

	if err = db.deleteAllUserSessions(tx, user.Name); err != nil {
		fmt.Println("database/DeleteAccount - fail deleting all user sessions")
		return
	}

	fmt.Println("database/DeleteAccount +")

	err = tx.Commit()
	return
}
