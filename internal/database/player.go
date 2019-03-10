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

func (db DataBase) checkBunch(field string, password string) (id int, err error) {

	fmt.Println("checkBunch:", field, password)

	// If checkBunchNamePass cant find brunch name-password
	if id, err = db.checkBunchNamePass(field, password); err != nil {
		// and checkBunchEmailPass cant find brunch email-password
		if id, err = db.checkBunchEmailPass(field, password); err != nil {
			return // then password wrong
		}
	}
	fmt.Println("i see id", id)
	err = nil
	return
}

// confirmRightPass checks that Player with such
// password and name exists
func (db DataBase) checkBunchNamePass(username string, password string) (id int, err error) {
	sqlStatement := "SELECT id FROM Player where name like $1 and password like $2"
	row := db.Db.QueryRow(sqlStatement, username, password)

	if err = row.Scan(&id); err != nil {
		err = errors.New("Wrong password")
	}
	fmt.Println("i found id", id)
	return
}

// confirmRightPass checks that Player with such
// password and name exists
func (db DataBase) checkBunchEmailPass(email string, password string) (id int, err error) {
	sqlStatement := "SELECT id FROM Player where email like $1 and password like $2"
	row := db.Db.QueryRow(sqlStatement, email, password)

	fmt.Println("email and password", email, password)
	if err := row.Scan(&id); err != nil {
		err = errors.New("Wrong password")
	}
	fmt.Println("i found id", id)
	return
}

// confirmRightEmail checks that Player with such
// email and name exists
func (db DataBase) confirmEmailNamePassword(user *models.UserPrivateInfo) error {
	sqlStatement := "SELECT 1 FROM Player where name like $1 and password like $2 and email like $3"

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

func (db *DataBase) createPlayer(user *models.UserPrivateInfo) (id int, err error) {
	sqlInsert := `
	INSERT INTO Player(name, password, email) VALUES
    ($1, $2, $3);
		`
	_, err = db.Db.Exec(sqlInsert, user.Name, user.Password, user.Email)

	if err != nil {
		return
	}

	sqlGetID := `SELECT id FROM Player WHERE name = $1`
	row := db.Db.QueryRow(sqlGetID, user.Name)

	err = row.Scan(&id)

	return
}
