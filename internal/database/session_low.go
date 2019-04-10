package database

import (
	"database/sql"
	misc "escapade/internal/misc"
)

func (db *DataBase) createSession(tx *sql.Tx, userID int) (sessionID string, err error) {
	sessionID = misc.CreateID()
	sqlStatement := `
	INSERT INTO Session(player_id, session_code)
		VALUES($1, $2);`

	_, err = tx.Exec(sqlStatement, userID, sessionID)
	return sessionID, err
}
