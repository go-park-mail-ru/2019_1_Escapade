package database

import (
	//
	"database/sql"
	"escapade/internal/models"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) createGamers(tx *sql.Tx, GameID int, gamers []*models.Gamer) (err error) {
	sqlInsert := `
	INSERT INTO Gamer(player_id, game_id, score, time,
		 left_click, right_click, explosion, won) VALUES
    ($1, $2, $3, $4::interval, $5, $6, $7, $8);
		`

	for _, gamer := range gamers {
		_, err = tx.Exec(sqlInsert, gamer.ID, GameID,
			gamer.LeftClick, gamer.RightClick, gamer.Explosion, gamer.Won)
		if err != nil {
			break
		}
	}
	return
}
