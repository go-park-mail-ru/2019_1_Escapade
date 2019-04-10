package database

import (
	"database/sql"
	"escapade/internal/models"
	"fmt"
)

// Register check sql-injections and are email and name unique
// Then add cookie to database and returns session_id
func (db *DataBase) Register(user *models.UserPrivateInfo) (str string, userID int, err error) {

	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if err = db.confirmUnique(tx, user); err != nil {
		fmt.Println("database/register - fail uniqie")
		return
	}

	if userID, err = db.createPlayer(tx, user); err != nil {
		fmt.Println("database/register - fail creating User")
		return
	}

	if str, err = db.createSession(tx, userID); err != nil {
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
func (db *DataBase) Login(user *models.UserPrivateInfo) (sessionCode string, found *models.UserPublicInfo, err error) {

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

	if sessionCode, err = db.createSession(tx, userID); err != nil {
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
func (db *DataBase) UpdatePlayerPersonalInfo(curName string, user *models.UserPrivateInfo) (err error) {
	var (
		curEmail string
		curPass  string
		oldName  string
		tx       *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	oldName = curName
	if curEmail, curPass, err = db.GetPasswordEmailByName(tx, curName); err != nil {
		return
	}

	if curEmail, err = db.checkParameter(tx, curEmail, user.Email, db.isEmailUnique); err != nil {
		return
	}

	if curName, err = db.checkParameter(tx, curName, user.Name, db.isNameUnique); err != nil {
		return
	}

	if user.Password != curPass && user.Password != "" {
		curPass = user.Password
	}

	user.Update(curName, curEmail, curPass)

	if _, err = db.updatePlayerPersonalInfo(tx, user, oldName); err != nil {
		return
	}

	err = tx.Commit()
	return
}

func (db *DataBase) GetUsers(difficult int, page int, perPage int,
	sort string) (players []models.UserPublicInfo, err error) {
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
