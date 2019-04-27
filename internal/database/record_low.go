package database

import (
<<<<<<< HEAD

=======
	"fmt"
>>>>>>> 508037185fc39abb3d6ee56a9fd2c48bac220f58
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

	var sqlStatement string
	record.Fix()
	if record.SingleWin > 0 {
		sqlStatement = `
		UPDATE Record
		SET score = (select 
			CASE 
				WHEN score>$1 THEN score 
				ELSE $1 
			END),
<<<<<<< HEAD
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
		RETURNING id`

		_, err = tx.Exec(sqlStatement, record.Score, record.Time,
			record.SingleTotal, record.OnlineTotal, record.SingleWin,
			record.OnlineWin, id, record.Difficult)
	} else {
		sqlStatement = `
		UPDATE Record
		SET singleTotal = singleTotal + $1,
				onlineTotal = onlineTotal + $2,
				singleWin = singleWin + $3,
				onlineWin = onlineWin + $4
		WHERE player_id = $5 and difficult = $6
		RETURNING id`

		_, err = tx.Exec(sqlStatement, record.SingleTotal,
			record.OnlineTotal, record.SingleWin,
			record.OnlineWin, id, record.Difficult)
	}
=======
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
>>>>>>> 508037185fc39abb3d6ee56a9fd2c48bac220f58

	return
}
