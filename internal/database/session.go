package database

import (
	misc "escapade/internal/misc"
	"fmt"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) createSession(userID int) (sessionID string, err error) {
	expiration := misc.CreateExpiration()
	sessionID = misc.CreateID()
	sqlStatement := `
	INSERT INTO Session(player_id, session_code, expiration)
		VALUES($1, $2, $3);
`
	fmt.Println("userID is ", userID)
	_, err = db.Db.Exec(sqlStatement, userID, sessionID, expiration)
	return sessionID, err
}

func (db *DataBase) deleteSession(sessionCode string) error {
	sqlStatement := `DELETE From Session where session_code=$1`
	_, err := db.Db.Exec(sqlStatement, sessionCode)
	return err
}

func (db *DataBase) deleteAllUserSessions(username string) (err error) {
	var id int
	if id, err = db.GetPlayerIDbyName(username); err != nil {
		return
	}

	sqlStatement := `DELETE From Session where player_id=$1`
	_, err = db.Db.Exec(sqlStatement, id)
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
