package database

import (
	"database/sql"
	"errors"
	"escapade/internal/models"
	"fmt"

	//
	_ "github.com/lib/pq"
)

// В будущем добавить, чтобы отдельно была проверка на
// на корректность, отдельно на sql  инъекции
func ValidatePrivateUI(user *models.UserPrivateInfo) (err error) {

	if !models.ValidatePassword(user.Password) {
		err = errors.New("password is not valid")
		return
	}

	if !models.ValidatePlayerName(user.Name) && !models.ValidateEmail(user.Email) {
		err = errors.New("player name or email is not valid")
		return
	}

	return
}

func isNameUnique(taken string, db *sql.DB) error {
	sqlStatement := "SELECT name " +
		"FROM Player where name=$1"

	row := db.QueryRow(sqlStatement, taken)

	var tmp string
	if err := row.Scan(&tmp); err != sql.ErrNoRows {
		if err == nil {
			return errors.New("name is taken")
		} else {
			return err
		}
	}
	return nil
}

func isEmailUnique(taken string, db *sql.DB) error {
	sqlStatement := "SELECT name " +
		"FROM Player where email=$1"

	row := db.QueryRow(sqlStatement, taken)

	var tmp string
	if err := row.Scan(&tmp); err != sql.ErrNoRows {
		if err == nil {
			return errors.New("email is taken")
		} else {
			return err
		}
	}
	return nil
}

// confirmUnique confirm that user.Email and user.Password unique
func confirmUnique(user *models.UserPrivateInfo, db *sql.DB) (err error) {
	err = isEmailUnique(user.Email, db)
	if err != nil {
		return
	}

	err = isNameUnique(user.Name, db)
	if err != nil {
		return
	}
	return
}

func confirmRightEmail(user *models.UserPrivateInfo, db *sql.DB) error {
	sqlStatement := "SELECT email " +
		"FROM Player where name=$1"

	// Get one record
	row := db.QueryRow(sqlStatement, user.Name)

	var email string

	if err := row.Scan(&email); err != nil {
		// No rows were returned
		return err
	}

	if email != user.Email {
		return errors.New("email is wrong")
	}

	return nil
}

func confirmRightPass(user *models.UserPrivateInfo, db *sql.DB) error {
	sqlStatement := "SELECT password " +
		"FROM Player where name=$1"

	// Get one record
	row := db.QueryRow(sqlStatement, user.Name)

	var password string

	if err := row.Scan(&password); err != nil {
		// No rows were returned
		return err
	}

	if password != user.Password {
		return errors.New("password is wrong")
	}

	return nil
}

func (db *DataBase) deleteUser(user *models.UserPrivateInfo) error {
	sqlStatement := `
	DELETE FROM Player where name=$1 and password=$2 and email=$3
		`
	_, err := db.Db.Exec(sqlStatement, user.Name,
		user.Password, user.Email)

	if err != nil {
		fmt.Println("database/user_private_info/deleteUser - fail:" + err.Error())

	}
	return err
}

func (db *DataBase) createUser(user *models.UserPrivateInfo) error {
	sqlStatement := `
	INSERT INTO Player(name, password, email) VALUES
    ($1, $2, $3);
		`
	_, err := db.Db.Exec(sqlStatement, user.Name,
		user.Password, user.Email)

	if err != nil {
		fmt.Println("database/user_private_info/createUser - fail:" + err.Error())

	}
	return err
}
