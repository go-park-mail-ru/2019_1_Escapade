package database

import (
	misc "escapade/internal/misc"
	"escapade/internal/models"
	"fmt"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) createSession(user *models.UserPrivateInfo) (string, error) {

	str := misc.CreateID()
	_, err := db.Db.Exec(`
				INSERT INTO Session(player_id, session_code, expiration)
				VALUES(
					(SELECT id FROM Player WHERE name=$1), $2, $3
				);
			`, user.Name, str, misc.CreateExpiration())

	if err != nil {
		fmt.Println("database/session/createSession - fail:" + err.Error())

	}
	return str, err
}

func (db *DataBase) deleteSession(sessionCode string) error {
	sqlStatement := `DELETE From Session where session_code=$1`
	_, err := db.Db.Exec(sqlStatement, sessionCode)
	return err
}
