package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/ery/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	//
	_ "github.com/jackc/pgx"
)

// CreateUser создать пользователя
func (db *DB) CreateUser(user *models.User) error {

	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = db.createUser(tx, user); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// UpdateUser обновить общую информацию о пользователе
func (db *DB) UpdateUser(user *models.User) error {

	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = db.updateUser(tx, user); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

// UpdateUserPrivate обновить пароль/имя пользователя
func (db *DB) UpdateUserPrivate(user *models.UpdatePrivateUser) error {

	// проверка на корректность пары логин/пароль
	if _, err := db.GetUserID(user.Old.Name, user.Old.Password); err != nil {
		return err
	}

	tx, err := db.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = db.updateUserPrivate(tx, &user.Old, &user.New); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}

func (db *DB) GetUser(userID int32) (models.User, error) {

	var user models.User
	tx, err := db.db.Beginx()
	if err != nil {
		return user, err
	}
	defer tx.Rollback()

	if user, err = db.getOneUser(tx, userID); err != nil {
		return user, err
	}

	err = tx.Commit()
	return user, err
}

func (db *DB) GetUserID(name, password string) (int32, error) {
	sqlStatement := "SELECT id " +
		" FROM Users where name like $1 and password like $2"

	row := db.db.QueryRow(sqlStatement, name, password)

	var userID int32
	err := row.Scan(&userID)
	if err != nil {
		utils.Debug(false, "cant get user's name and password")
	}

	return userID, err
}

func (db *DB) GetUsers(name string) (models.Users, error) {

	var users models.Users
	tx, err := db.db.Beginx()
	if err != nil {
		return users, err
	}
	defer tx.Rollback()

	if name == "" {
		users.Users, err = db.getAllUsers(tx)
	} else {
		users.Users, err = db.searchUsersWithName(tx, name)
	}
	if err != nil {
		return users, err
	}

	err = tx.Commit()
	return users, err
}
