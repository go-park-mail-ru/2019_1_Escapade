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

// GetPlayerIDbyName get player's id by his hame
func (db *DataBase) GetPlayerIDbyName(username string) (id int, err error) {
	sqlStatement := `SELECT id FROM Player WHERE name = $1`
	row := db.Db.QueryRow(sqlStatement, username)

	err = row.Scan(&id)
	return
}

// GetNameByEmail get player's name by his email
func (db DataBase) GetNameByEmail(email string) (name string, err error) {
	sqlStatement := "SELECT name " +
		"FROM Player where email=$1"

	row := db.Db.QueryRow(sqlStatement, email)

	if err = row.Scan(&name); err != nil {
		return
	}
	return
}

// isNameUnique checks if there are Players with
// this('taken') name and returns corresponding error if yes
func (db DataBase) isNameUnique(taken string) error {
	sqlStatement := "SELECT name " +
		"FROM Player where name=$1"

	row := db.Db.QueryRow(sqlStatement, taken)

	var tmp string
	if err := row.Scan(&tmp); err != sql.ErrNoRows {
		if err == nil {
			return errors.New("name is taken")
		}
		return err
	}
	return nil
}

// isEmailUnique checks if there are Players with
// this('taken') email and returns corresponding error if yes
func (db DataBase) isEmailUnique(taken string) error {
	sqlStatement := "SELECT name " +
		"FROM Player where email=$1"

	row := db.Db.QueryRow(sqlStatement, taken)

	var tmp string
	if err := row.Scan(&tmp); err != sql.ErrNoRows {
		if err == nil {
			return errors.New("email is taken")
		}
		return err
	}
	return nil
}

// confirmUnique confirm that user.Email and user.Name
// dont use by another Player
func (db DataBase) confirmUnique(user *models.UserPrivateInfo) (err error) {
	if err = db.isEmailUnique(user.Email); err != nil {
		return
	}

	err = db.isNameUnique(user.Name)
	return
}

func (db DataBase) checkBunch(field string, password string) (err error) {
	var (
		right1 bool
		right2 bool
	)
	fmt.Printf("checkBunch:", field, password)

	if right1, err = db.checkBunchNamePass(field, password); err != nil {
		return
	}

	if right2, err = db.checkBunchEmailPass(field, password); err != nil {
		return
	}

	if !right1 && !right2 {
		return errors.New("Wrong password")
	}
	return
}

// confirmRightPass checks that Player with such
// password and name exists
func (db DataBase) checkBunchNamePass(username string, password string) (bool, error) {
	sqlStatement := "SELECT password " +
		"FROM Player where name=$1"

	row := db.Db.QueryRow(sqlStatement, username)

	var get string

	if err := row.Scan(&get); err != nil || password != get {
		return false, err
	}

	return true, nil
}

// confirmRightPass checks that Player with such
// password and name exists
func (db DataBase) checkBunchEmailPass(email string, password string) (bool, error) {
	sqlStatement := "SELECT password " +
		"FROM Player where email=$1"

	row := db.Db.QueryRow(sqlStatement, email)

	var get string

	if err := row.Scan(&get); err != nil || password != get {
		return false, err
	}

	return true, nil
}

// confirmRightEmail checks that Player with such
// email and name exists
func (db DataBase) confirmEmailNamePassword(user *models.UserPrivateInfo) error {
	sqlStatement := "SELECT 1 FROM Player where name=$1, password=$2, email=$3"

	row := db.Db.QueryRow(sqlStatement, user.Name, user.Password, user.Email)
	var res int
	err := row.Scan(&res)
	return err
}

func (db *DataBase) deletePlayer(user *models.UserPrivateInfo) error {
	sqlStatement := `
	DELETE FROM Player where name=$1 and password=$2 and email=$3
		`
	_, err := db.Db.Exec(sqlStatement, user.Name,
		user.Password, user.Email)

	return err
}

func (db *DataBase) createPlayer(user *models.UserPrivateInfo) error {
	sqlStatement := `
	INSERT INTO Player(name, password, email) VALUES
    ($1, $2, $3);
		`
	_, err := db.Db.Exec(sqlStatement, user.Name,
		user.Password, user.Email)

	return err
}
