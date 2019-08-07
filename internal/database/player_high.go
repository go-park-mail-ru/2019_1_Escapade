package database

import (
	"fmt"
	"math"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"

	"database/sql"
)

// Register check sql-injections and is name unique
// Then add cookie to database and returns session_id
func (db *DataBase) Register(user *models.UserPrivateInfo) (userID int, err error) {

	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	fmt.Println("lets create", user.Name)
	if userID, err = db.createPlayer(tx, user); err != nil {
		err = re.ErrorUserIsExist()
		return
	}

	if err = db.createRecords(tx, userID); err != nil {
		return
	}

	err = tx.Commit()
	return
}

// Login check sql-injections and is password right
// Then add cookie to database and returns session_id
func (db *DataBase) Login(name, password string) (int32, error) {

	var (
		tx     *sql.Tx
		userID int32
		err    error
	)

	if tx, err = db.Db.Begin(); err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if userID, _, err = db.checkBunch(tx, name, password); err != nil {
		return userID, err
	}

	err = tx.Commit()
	return userID, err
}

// UpdatePlayerPersonalInfo gets name of Player from
// relation Session, cause we know that user has session
func (db *DataBase) UpdatePlayerPersonalInfo(userID int32, user *models.UserPrivateInfo) (err error) {
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

// GetUsers get users
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

// GetUser get one user
func (db *DataBase) GetUser(userID int32, difficult int) (user *models.UserPublicInfo, err error) {

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
		return
	}

	// TODO delete all tokens
	// if err = db.deleteAllUserSessions(tx, user.Name); err != nil {
	// 	return
	// }

	err = tx.Commit()
	return
}

// GetUsersPageAmount returns amount of rows in table Player
// deleted on amount of rows in one page
func (db *DataBase) GetUsersPageAmount(perPage int) (amount int, err error) {
	sqlStatement := `SELECT count(1) FROM Player`
	row := db.Db.QueryRow(sqlStatement)
	if err = row.Scan(&amount); err != nil {
		return
	}

	if amount > db.PageUsers {
		amount = db.PageUsers
	}
	if perPage == 0 {
		perPage = 1
	}
	amount = int(math.Ceil(float64(amount) / float64(perPage)))
	return
}
