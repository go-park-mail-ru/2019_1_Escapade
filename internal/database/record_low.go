package database

import (
	"fmt"
	//
	"database/sql"
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) createRecords(tx *sql.Tx, id int) (err error) {
	sqlInsert := `
	INSERT INTO Record(player_id, difficult) VALUES
    ($1, $2);
		`
	difficultAmount := 4 // вынести в конфиг
	for i := 0; i < difficultAmount; i++ {
		_, err = tx.Exec(sqlInsert, id, i)
		if err != nil {
			return
		}
	}
	return
}

func (db *DataBase) updateRecords(tx *sql.Tx, id int,
	record *models.Record) (err error) {

	sqlStatement := `
	UPDATE Record
	SET score = (select 
		CASE 
			WHEN score>$1 THEN score 
			ELSE $1 
		END),
		 time = (select 
			CASE 
				WHEN time<$2::interval THEN time 
				ELSE $2::interval
			END),
			singleTotal = singleTotal + $3,
			onlineTotal = onlineTotal + $4,
			singleWin = singleWin + $5,
			onlineWin = onlineWin + $6
	WHERE player_id = $7 and difficult = $8
	RETURNING id
`
	record.Fix()
	fmt.Println("record.Score = ", record.Score)
	_, err = tx.Exec(sqlStatement, record.Score, record.Time,
		record.SingleTotal, record.OnlineTotal, record.SingleWin,
		record.OnlineWin, id, record.Difficult)

	return
}
