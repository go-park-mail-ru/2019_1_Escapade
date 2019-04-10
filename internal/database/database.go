package database

import (
	"database/sql"
	"escapade/internal/models"
	"math"

	"fmt"

	//
	_ "github.com/lib/pq"
)

// DataBase consists of *sql.DB
// Support methods Login, Register
type DataBase struct {
	Db        *sql.DB
	PageGames int
	PageUsers int
}

// Logout delete session_id row  from session table
func (db *DataBase) Logout(sessionCode string) (err error) {
	err = db.deleteSession(sessionCode)
	return
}

// PostImage set filename of avatar to relation Player
func (db *DataBase) PostImage(filename string, userID int) (err error) {
	sqlStatement := `UPDATE Player SET photo_title = $1 WHERE id = $2;`

	_, err = db.Db.Exec(sqlStatement, filename, userID)

	if err != nil {
		fmt.Println("database/session/PostImage - fail:" + err.Error())
		return
	}
	return
}

// GetImage Get avatar - filename of player image
func (db *DataBase) GetImage(userID int) (filename string, err error) {
	sqlStatement := `
	SELECT photo_title
		FROM Player as P 
		WHERE P.id = $1 
`
	row := db.Db.QueryRow(sqlStatement, userID)

	if err = row.Scan(&filename); err != nil {
		fmt.Println("database/GetImage failed")
		return
	}
	return
}

// GetNameBySessionID gets name of Player from
// relation Session, cause we know that user has session
func (db *DataBase) GetNameBySessionID(sessionID string) (name string, err error) {
	sqlStatement := `
	SELECT name
	FROM Player as P join Session as S on S.player_id=P.id
	WHERE session_code like $1 
`
	row := db.Db.QueryRow(sqlStatement, sessionID)

	err = row.Scan(&name)
	if err != nil {
		fmt.Println("Sess error: ", err.Error())
		fmt.Println("database/GetNameBySessionID failed")
		return
	}

	return
}

// GetNameBySessionID gets name of Player from
// relation Session, cause we know that user has session
func (db *DataBase) GetUserIdBySessionID(sessionID string) (id int, err error) {
	sqlStatement := `
	SELECT S.player_id
	FROM Session as S
	WHERE session_code like $1 
	`
	row := db.Db.QueryRow(sqlStatement, sessionID)

	err = row.Scan(&id)
	if err != nil {
		fmt.Println("Sess error: ", err.Error())
		fmt.Println("database/GetIdBySessionID failed")
		return
	}

	return
}

// GetUsersPageAmount returns amount of rows in table Player
// deleted on amount of rows in one page
func (db *DataBase) GetUsersPageAmount(per_page int) (amount int, err error) {
	sqlStatement := `SELECT count(1) FROM Player`
	row := db.Db.QueryRow(sqlStatement)
	if err = row.Scan(&amount); err != nil {
		fmt.Println("GetUsersAmount failed")
		return
	}
	if amount > db.PageUsers {
		amount = db.PageUsers
	}
	amount = int(math.Ceil(float64(amount) / float64(per_page)))
	return
}

func (db *DataBase) GetProfile(name string) (user models.UserPublicInfo, err error) {

	sqlStatement := `
	SELECT email, best_score, best_time 
	FROM Player 
	WHERE name like $1 
`

	row := db.Db.QueryRow(sqlStatement, name)

	user.Name = name

	if err = row.Scan(&user.Email, &user.BestScore,
		&user.BestTime); err != nil {
		fmt.Println("database/GetProfile failed")
		return
	}

	fmt.Println("database/GetProfile +")

	return
}

// GetFullGamesInformation returns games, played by player with some name
func (db *DataBase) GetFullGamesInformation(name string,
	page int) (games []models.GameInformation, err error) {

	size := db.PageGames
	sqlStatement := `
	SELECT 	ge.width, ge.height, ge.difficult,
					ge.players, ge.mines, ge.date, ge.online,
					gr.score, gr.time, gr.mines_open,
					gr.left_click, gr.right_click,
					gr.explosion, gr.won
	 FROM Player as p 
		JOIN
			(
				SELECT * FROM Game
			) as ge
		ON p.id = ge.player_id and p.name like $1
		JOIN
			(
				SELECT * FROM Gamer
			) as gr
			ON p.id = gr.player_id and ge.id = gr.game_id
			OFFSET $2 Limit $3
	`

	games = make([]models.GameInformation, 0, size)
	rows, erro := db.Db.Query(sqlStatement, name, size*(page-1), size) // //, name)

	if erro != nil {
		err = erro

		fmt.Println("database/GetGames cant access to database:", erro.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		game := models.GameInformation{}
		game.Game = &models.Game{}
		game.Gamer = &models.Gamer{}
		/*
			ge.width, ge.height, ge.difficult,
							ge.players, ge.mines, ge.date, ge.online,
							gr.score, gr.time, gr.mines_open,
							gr.left_click, gr.right_click,
							gr.explosion, gr.won
			 FROM Player as p */
		if err = rows.Scan(&game.Game.Height,
			&game.Game.Difficult, &game.Game.Players,
			&game.Game.Mines, &game.Game.Date, &game.Game.Online,
			&game.Gamer.Score, &game.Gamer.Time, &game.Gamer.MinesOpen,
			&game.Gamer.LeftClick, &game.Gamer.RightClick,
			&game.Gamer.Explosion, &game.Gamer.Won); err != nil {

			fmt.Println("database/GetGames wrong row catched")

			break
		}

		games = append(games, game)
	}

	fmt.Println("database/GetGames +")

	return
}

// DeleteAccount deletes account
func (db *DataBase) DeleteAccount(user *models.UserPrivateInfo) (err error) {

	if err = db.confirmEmailNamePassword(user); err != nil {
		fmt.Println("database/DeleteAccount - fail confirmition")
		return
	}

	if err = db.deleteAllUserSessions(user.Name); err != nil {
		fmt.Println("database/DeleteAccount - fail deleting all user sessions")
		return
	}

	if err = db.deletePlayer(user); err != nil {
		fmt.Println("database/DeleteAccount - fail deletting User")
		return
	}

	fmt.Println("database/DeleteAccount +")

	return
}
