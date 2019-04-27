package database

import (
	"database/sql"
	"escapade/internal/models"
	"fmt"
	"time"
)

// createPlayer create player
func (db *DataBase) createPlayer(tx *sql.Tx, user *models.UserPrivateInfo) (id int, err error) {
	sqlInsert := `
	INSERT INTO Player(name, password, email, firstSeen, lastSeen) VALUES
    ($1, $2, $3, $4, $5);
		`
	t := time.Now()
	_, err = db.Db.Exec(sqlInsert, user.Name,
		user.Password, user.Email, t, t)

	if err != nil {
		return
	}

	sqlGetID := `SELECT id FROM Player WHERE name = $1`
	row := tx.QueryRow(sqlGetID, user.Name)

	err = row.Scan(&id)
	return
}

func (db *DataBase) updatePlayerPersonalInfo(tx *sql.Tx, user *models.UserPrivateInfo) (err error) {
	sqlStatement := `
			UPDATE Player 
			SET name = $1, email = $2, password = $3, lastSeen = $4
			WHERE id = $5
			RETURNING id
		`

	_, err = tx.Exec(sqlStatement, user.Name,
		user.Email, user.Password, time.Now(), user.ID)
	return
}

// updatePlayerLastSeen update users last date seen
func (db *DataBase) updatePlayerLastSeen(tx *sql.Tx, id int) (err error) {
	sqlStatement := `
			UPDATE Player 
			SET lastSeen = $1
			WHERE id = $2
		`

	_, err = tx.Exec(sqlStatement, time.Now(), id)
	return
}

func (db DataBase) checkBunch(tx *sql.Tx, field string, password string) (id int, user *models.UserPublicInfo, err error) {
	sqlStatement := `
	SELECT pl.id, pl.name, pl.email, r.score, r.time, r.difficult
		FROM Player as pl
		join Record as r 
		on r.player_id = pl.id
		where r.difficult = 0 and password like $1 and (name like $2 or email like $2)`

	row := tx.QueryRow(sqlStatement, password, field)

	user = &models.UserPublicInfo{}
	err = row.Scan(&id, &user.Name, &user.Email, &user.BestScore,
		&user.BestTime, &user.Difficult)
	return
}

// GetPrivateInfo get player's personal info
func (db DataBase) getPrivateInfo(tx *sql.Tx, userID int) (user *models.UserPrivateInfo, err error) {
	sqlStatement := "SELECT name, email, password " +
		"FROM Player where id = $1"

	row := tx.QueryRow(sqlStatement, userID)

	user = &models.UserPrivateInfo{}
	user.ID = userID
	err = row.Scan(&user.Name, &user.Email, &user.Password)

	return
}

// GetUsers returns information about users
// for leaderboard
func (db *DataBase) getUsers(tx *sql.Tx, difficult int, offset int, limit int,
	sort string) (players []*models.UserPublicInfo, err error) {

	sqlStatement := `
	SELECT P.id, P.photo_title, P.name, P.email,
				 R.score, R.time, R.Difficult
	FROM Player as P
	join Record as R 
	on R.player_id = P.id
	where r.difficult = $1  
	`
	if sort == "score" {
		sqlStatement += ` ORDER BY (score) desc `
	} else {
		sqlStatement += ` ORDER BY (time) `
	}
	sqlStatement += ` OFFSET $2 Limit $3 `

	players = make([]*models.UserPublicInfo, 0, limit)
	rows, erro := tx.Query(sqlStatement, difficult, offset, limit)

	if erro != nil {
		err = erro

		fmt.Println("database/GetUsers cant access to database:", erro.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		player := &models.UserPublicInfo{}
		if err = rows.Scan(&player.ID, &player.FileKey, &player.Name, &player.Email, &player.BestScore,
			&player.BestTime, &player.Difficult); err != nil {

			fmt.Println("database/GetUsers wrong row catched")

			break
		}

		fmt.Println("wrote: ", player.BestScore, player.BestTime)

		players = append(players, player)
	}

	fmt.Println("database/GetUsers +")

	return
}

// GetUsers returns information about users
// for leaderboard
func (db *DataBase) getUser(tx *sql.Tx, userID int, difficult int) (player *models.UserPublicInfo, err error) {

	sqlStatement := `
	SELECT P.id, P.photo_title, P.name, P.email,
				 R.score, R.time, R.Difficult
	FROM Player as P
	join Record as R 
	on R.player_id = P.id
	where R.player_id = $1 and
		R.difficult = $2
	`

	player = &models.UserPublicInfo{}
	row := tx.QueryRow(sqlStatement, userID, difficult)
	err = row.Scan(&player.ID, &player.FileKey, &player.Name,
		&player.Email, &player.BestScore, &player.BestTime, &player.Difficult)
	return
}
<<<<<<< HEAD

func (db *DataBase) deletePlayer(tx *sql.Tx, user *models.UserPrivateInfo) error {
	sqlStatement := `
	DELETE FROM Player where name=$1 and password=$2 and email=$3
		`
	fmt.Println("+++++")

	fmt.Println("user.Name, user.Password, user.Email", user.Name, user.Password, user.Email)

	_, err := tx.Exec(sqlStatement, user.Name,
		user.Password, user.Email)

	return err
}
=======
>>>>>>> 508037185fc39abb3d6ee56a9fd2c48bac220f58
