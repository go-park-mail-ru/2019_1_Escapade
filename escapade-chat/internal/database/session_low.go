package database

import (
	"database/sql"
)

func (db *DataBase) createSession(tx *sql.Tx, userID int, sessionID string) (err error) {
	sqlStatement := `
	INSERT INTO Session(player_id, session_code)
		VALUES($1, $2);`

	_, err = tx.Exec(sqlStatement, userID, sessionID)
	return
}
