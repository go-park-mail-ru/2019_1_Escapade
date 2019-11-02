package database

// ImageRepositoryPQ implements the interface ImageRepositoryI using the sql postgres driver
type ImageRepositoryPQ struct{}

// update set filename of avatar to relation Player
func (db *ImageRepositoryPQ) update(dbI DatabaseI, filename string, userID int32) error {
	sqlStatement := `UPDATE Player SET photo_title = $1 WHERE id = $2;`

	var err error
	_, err = dbI.Exec(sqlStatement, filename, userID)

	return err
}

// fetchByName Get avatar - filename of player image by his name
func (db *ImageRepositoryPQ) fetchByName(dbI DatabaseI, name string) (string, error) {
	sqlStatement := `SELECT photo_title FROM Player WHERE name like $1`
	row := dbI.QueryRow(sqlStatement, name)

	var filename string
	err := row.Scan(&filename)
	return filename, err
}

// fetchByID Get avatar - filename of player image by his id
func (db *ImageRepositoryPQ) fetchByID(dbI DatabaseI, id int32) (string, error) {
	sqlStatement := `SELECT photo_title FROM Player WHERE id=$1`
	row := dbI.QueryRow(sqlStatement, id)

	var filename string
	err := row.Scan(&filename)
	return filename, err
}
