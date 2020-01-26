package database

import (
	"context"
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
)

// Record implements the interface RecordRepositoryI using the sql postgres driver
type Record struct {
	db idb.ExecerI
}

func NewRecord(dbI idb.ExecerI) *Record {
	return &Record{dbI}
}

// Create user's record
func (db *Record) Create(ctx context.Context, id int) error {
	var err error
	sqlInsert := `INSERT INTO Record(player_id, difficult) VALUES ($1, $2);`
	difficultAmount := 4 // вынести в конфиг
	for i := 0; i < difficultAmount; i++ {
		_, err = db.db.ExecContext(ctx, sqlInsert, id, i)
		if err != nil {
			break
		}
	}
	return err
}

// Update user's record
func (db *Record) Update(ctx context.Context, id int32, record *models.Record) error {
	var (
		sqlStatement string
		err          error
	)
	record.Fix()
	if record.SingleWin > 0 {
		sqlStatement = `
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
		RETURNING id`

		_, err = db.db.ExecContext(ctx, sqlStatement, record.Score, record.Time,
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

		_, err = db.db.ExecContext(ctx, sqlStatement, record.SingleTotal,
			record.OnlineTotal, record.SingleWin,
			record.OnlineWin, id, record.Difficult)
	}

	return err
}
