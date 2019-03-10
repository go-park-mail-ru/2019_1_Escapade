package database

import (
	misc "escapade/internal/misc"
	"escapade/internal/models"
	"fmt"
	"time"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) createSession(user *models.UserPrivateInfo) (string, error) {

	var (
		err          error
		sessionID    string
		sqlStatement string
		expiration   time.Time
	)
	expiration = misc.CreateExpiration()
	sessionID = misc.CreateID()
	sqlStatement = `
	INSERT INTO Session(player_id, session_code, expiration)
	VALUES(
		(SELECT id FROM Player WHERE name=$1), $2, $3
	);
`

	_, err = db.Db.Exec(sqlStatement, user.Name, sessionID, expiration)

	if err != nil {
		fmt.Println("database/session/createSession - fail:" + err.Error())
		return sessionID, err
	}
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
