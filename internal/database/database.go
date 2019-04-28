package database

import (
	"database/sql"
	"math"

	"fmt"

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
		fmt.Println("database/session/PostImage - fail:" + err.Error())
		return
	}
	return
}

// GetImage Get avatar - filename of player image
func (db *DataBase) GetImage(userID int) (filename string, err error) {
	sqlStatement := `
	SELECT photo_title
		FROM Player as P 
		WHERE P.id = $1 
`
	row := db.Db.QueryRow(sqlStatement, userID)

	if err = row.Scan(&filename); err != nil {
		fmt.Println("database/GetImage failed")
		return
	}
	return
}

// GetNameBySessionID gets name of Player from
// relation Session, cause we know that user has session
func (db *DataBase) GetNameBySessionID(sessionID string) (name string, err error) {
	sqlStatement := `
	SELECT name
	FROM Player as P join Session as S on S.player_id=P.id
	WHERE session_code like $1 
`
	row := db.Db.QueryRow(sqlStatement, sessionID)

	err = row.Scan(&name)
	if err != nil {
		fmt.Println("Sess error: ", err.Error())
		fmt.Println("database/GetNameBySessionID failed")
		return
	}

	return
}

// GetNameBySessionID gets name of Player from
// relation Session, cause we know that user has session
func (db *DataBase) GetUserIdBySessionID(sessionID string) (id int, err error) {
	sqlStatement := `
	SELECT S.player_id
	FROM Session as S
	WHERE session_code like $1 
	`
	row := db.Db.QueryRow(sqlStatement, sessionID)

	err = row.Scan(&id)
	if err != nil {
		fmt.Println("Sess error: ", err.Error())
		fmt.Println("database/GetIdBySessionID failed")
		return
	}

	return
}

// GetNameBySessionID gets name of Player from
// relation Session, cause we know that user has session
func (db *DataBase) GetUserIdByName(name string) (id int, err error) {
	sqlStatement := `
	SELECT id
	FROM Player
	WHERE Name like $1 
	`
	row := db.Db.QueryRow(sqlStatement, name)

	err = row.Scan(&id)
	if err != nil {
		fmt.Println("Sess error: ", err.Error())
		fmt.Println("database/GetIdBySessionID failed")
		return
	}

	return
}

// GetUsersPageAmount returns amount of rows in table Player
// deleted on amount of rows in one page
func (db *DataBase) GetUsersPageAmount(per_page int) (amount int, err error) {
	sqlStatement := `SELECT count(1) FROM Player`
	row := db.Db.QueryRow(sqlStatement)
	if err = row.Scan(&amount); err != nil {
		return
	}

	if amount > db.PageUsers {
		amount = db.PageUsers
	}
	if per_page == 0 {
		per_page = 1
	}
	amount = int(math.Ceil(float64(amount) / float64(per_page)))
	return
}
