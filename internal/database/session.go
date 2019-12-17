package database

// delete sessions -> introduced Redis

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func (db *DataBase) DeleteSession(sessionCode string) error {
	sqlStatement := `DELETE From Session where session_code=$1`
	_, err := db.Db.Exec(sqlStatement, sessionCode)
	return err
}

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
		where p.name like $1
	`
	row := db.Db.QueryRow(sqlStatement, userName)

	err = row.Scan(&sessionID)
	return sessionID, err
}

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
