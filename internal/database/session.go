package database

import (
	//
	"database/sql"

	_ "github.com/lib/pq"
)

func (db *DataBase) deleteAllUserSessions(tx *sql.Tx, username string) (err error) {
	var id int
	if id, err = db.GetPlayerIDbyName(username); err != nil {
		return
	}

	sqlStatement := `DELETE From Session where player_id=$1`
	_, err = tx.Exec(sqlStatement, id)
	return
}

func (db *DataBase) GetSessionByName(userName string) (sessionID string, err error) {

	sqlStatement := `
	select s.session_code 
		from Session as s join Player as p
		on s.player_id = p.id 
		where p.name like $1 or email like $1
	`
	row := db.Db.QueryRow(sqlStatement, userName)

	err = row.Scan(&sessionID)
	return sessionID, err
}
