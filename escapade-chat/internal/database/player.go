package database

import (
	//
	"fmt"

	_ "github.com/lib/pq"
)

// GetPlayerIDbyName get player's id by his hame
func (db *DataBase) GetPlayerIDbyName(username string) (id int, err error) {
	sqlStatement := `SELECT id FROM Player WHERE name = $1`
	row := db.Db.QueryRow(sqlStatement, username)

	err = row.Scan(&id)
	return
}

// TODO delete it, when all tests will be done
// GetPlayerNames get all players name
func (db *DataBase) GetPlayerNames() (names []string, err error) {
	sqlStatement := `SELECT name FROM Player`
	names = make([]string, 0)
	rows, erro := db.Db.Query(sqlStatement)

	if erro != nil {
		err = erro
		return
	}
	defer rows.Close()

	for rows.Next() {
		var str string
		if err = rows.Scan(&str); err != nil {
			break
		}

		fmt.Println("add:", str)
		names = append(names, str)
	}

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
