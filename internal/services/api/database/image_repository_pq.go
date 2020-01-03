package database

import idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"

// ImageRepositoryPQ implements the interface ImageRepositoryI using the sql postgres driver
type ImageRepositoryPQ struct{}

// Update set filename of avatar to relation Player
func (db *ImageRepositoryPQ) Update(dbI idb.Interface, filename string, userID int32) error {
	sqlStatement := `UPDATE Player SET photo_title = $1 WHERE id = $2;`
	_, err := dbI.Exec(sqlStatement, filename, userID)
	return err
}

// FetchByName Get avatar - filename of player image by his name
func (db *ImageRepositoryPQ) FetchByName(dbI idb.Interface, name string) (string, error) {
	sqlStatement := `SELECT photo_title FROM Player WHERE name like $1`
	var filename string
	err := dbI.QueryRow(sqlStatement, name).Scan(&filename)
	return filename, err
}

// FetchByID Get avatar - filename of player image by his id
func (db *ImageRepositoryPQ) FetchByID(dbI idb.Interface, id int32) (string, error) {
	sqlStatement := `SELECT photo_title FROM Player WHERE id=$1`
	var filename string
	err := dbI.QueryRow(sqlStatement, id).Scan(&filename)
	return filename, err
}

// 33
