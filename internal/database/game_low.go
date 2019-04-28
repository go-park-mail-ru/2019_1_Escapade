package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"database/sql"
	"fmt"
	"time"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) createGame(tx *sql.Tx, UserID int, game *models.Game) (err error) {
	sqlInsert := `
	INSERT INTO Game(player_id, difficult, width, height,
		 players, mines, date) VALUES
    ($1, $2, $3, $4, $5, $6, $7);
		`
	_, err = tx.Exec(sqlInsert, UserID, game.Difficult,
		game.Width, game.Height, game.Players, game.Mines,
		time.Now())

	return
}

// GetFullGamesInformation returns games, played by player with some name
func (db *DataBase) GetFullGamesInformation(tx *sql.Tx, UserID int,
	page int) (games []*models.GameInformation, err error) {

	size := db.PageGames
	sqlStatement := `
	SELECT 	GE.width, GE.height, GE.difficult, GE.players,
	 GE.mines, GE.date, GR.score, 
		GR.time, GR.left_click, GR.right_click, GR.explosion,
		GR.won
	 FROM Gamer as GR
	 join Game as GE
		ON GR.id = $1 and GR.game_id = GE.id
		OFFSET $2 Limit $3
	`

	games = make([]*models.GameInformation, 0, size)
	rows, erro := tx.Query(sqlStatement, UserID, size*(page-1), size) // //, name)

	if erro != nil {
		err = erro

		fmt.Println("database/GetGames cant access to database:", erro.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		game := &models.GameInformation{}
		game.Game = &models.Game{}
		game.Gamers = make([]*models.Gamer, 1)
		/*
			ge.width, ge.height, ge.difficult,
							ge.players, ge.mines, ge.date, ge.online,
							gr.score, gr.time, gr.mines_open,
							gr.left_click, gr.right_click,
							gr.explosion, gr.won
			 FROM Player as p */

		if err = rows.Scan(&game.Game.Width,
			&game.Game.Height,
			&game.Game.Difficult, &game.Game.Players,
			&game.Game.Mines, &game.Game.Date,
			&game.Gamers[0].Score, &game.Gamers[0].Time,
			&game.Gamers[0].LeftClick, &game.Gamers[0].RightClick,
			&game.Gamers[0].Explosion, &game.Gamers[0].Won); err != nil {

			fmt.Println("database/GetGames wrong row catched")

			break
		}

		games = append(games, game)
	}

	fmt.Println("database/GetGames +")

	return
}
