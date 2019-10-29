package database

import (
	"fmt"
)

// PostImage set filename of avatar to relation Player
func (db *DataBase) PostImage(filename string, userID int32) error {
	sqlStatement := `UPDATE Player SET photo_title = $1 WHERE id = $2;`

	var err error
	_, err = db.Db.Exec(sqlStatement, filename, userID)

	return err
}

// GetImageByName Get avatar - filename of player image by his name
func (db *DataBase) GetImageByName(name string) (string, error) {
	fmt.Println("name is ", name)
	sqlStatement := `SELECT photo_title FROM Player WHERE name like $1`
	row := db.Db.QueryRow(sqlStatement, name)

	var filename string
	err := row.Scan(&filename)

	if err != nil {
		fmt.Println("GetImageByName err", err.Error())
		return filename, err
	}
	fmt.Println("GetImageByName")
	return filename, err
}

// GetImageByID Get avatar - filename of player image by his id
func (db *DataBase) GetImageByID(id int32) (string, error) {
	fmt.Println("id is ", id)
	sqlStatement := `SELECT photo_title FROM Player WHERE id=$1`
	row := db.Db.QueryRow(sqlStatement, id)

	var filename string
	err := row.Scan(&filename)

	if err != nil {
		fmt.Println("GetImageByID err", err.Error())
		return filename, err
	}
	fmt.Println("GetImageByID")
	return filename, err
}
