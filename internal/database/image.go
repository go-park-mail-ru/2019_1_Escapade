package database

import (
	"database/sql"

	//
	_ "github.com/lib/pq"
)

// DataBase consists of *sql.DB
// Support methods Login, Register
type DataBase struct {
	Db        *sql.DB
	PageGames int
	PageUsers int
}

// Logout delete session_id row  from session table
// func (db *DataBase) Logout(sessionCode string) (err error) {
// 	err = db.deleteSession(sessionCode)
// 	return
// }

// PostImage set filename of avatar to relation Player
func (db *DataBase) PostImage(filename string, userID int) (err error) {
	sqlStatement := `UPDATE Player SET photo_title = $1 WHERE id = $2;`

	_, err = db.Db.Exec(sqlStatement, filename, userID)

	if err != nil {
		return
	}
	return
}

// GetImage Get avatar - filename of player image
func (db *DataBase) GetImage(name string) (filename string, err error) {
	sqlStatement := `
	SELECT photo_title
		FROM Player as P 
		WHERE P.name like $1 
`
	row := db.Db.QueryRow(sqlStatement, name)

	if err = row.Scan(&filename); err != nil {
		return
	}
	return
}
