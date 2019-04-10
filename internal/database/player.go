package database

import (
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

// GetPlayerIDbyName get player's id by his hame
func (db *DataBase) GetPlayerIDbyName(username string) (id int, err error) {
	sqlStatement := `SELECT id FROM Player WHERE name = $1`
	row := db.Db.QueryRow(sqlStatement, username)

	err = row.Scan(&id)
	return
}

// GetPlayerNamebyID get player's name by his id
func (db *DataBase) GetPlayerNamebyID(id int) (username string, err error) {
	sqlStatement := `SELECT name FROM Player WHERE id = $1`
	row := db.Db.QueryRow(sqlStatement, id)

	err = row.Scan(&username)
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
