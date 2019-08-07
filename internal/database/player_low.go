package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	re "github.com/go-park-mail-ru/2019_1_Escapade/internal/return_errors"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"database/sql"
	"fmt"
	"time"
)

// createPlayer create player
func (db *DataBase) createPlayer(tx *sql.Tx, user *models.UserPrivateInfo) (id int, err error) {
	sqlInsert := `
	INSERT INTO Player(name, password, firstSeen, lastSeen) VALUES
		($1, $2, $3, $4)
		RETURNING id;
		`
	t := time.Now()
	row := tx.QueryRow(sqlInsert, user.Name,
		user.Password, t, t)

	err = row.Scan(&id)
	if err == nil {
		fmt.Println("register:", user.Name, user.Password, id)
	} else {
		fmt.Println("create err", err.Error())
	}
	return
}

func (db *DataBase) updatePlayerPersonalInfo(tx *sql.Tx, user *models.UserPrivateInfo) (err error) {
	sqlStatement := `
			UPDATE Player 
			SET name = $1, password = $2, lastSeen = $3
			WHERE id = $4
			RETURNING id
		`

	fmt.Println("Update to", user.Name, user.Password)
	row := tx.QueryRow(sqlStatement, user.Name,
		user.Password, time.Now(), user.ID)
	err = row.Scan(&user.ID)
	if err != nil {
		fmt.Println("updatePlayerPersonalInfo: err", err.Error())
		err = re.ErrorUserIsExist()
	} else {
		fmt.Println("updatePlayerPersonalInfo done", user.Name,
			user.Password, user.ID)
	}

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
	if err != nil {
		utils.Debug(true, "cant update players's last seen")
	}
	return
}

func (db DataBase) checkBunch(tx *sql.Tx, field string, password string) (id int32, user *models.UserPublicInfo, err error) {
	sqlStatement := `
	SELECT pl.id, pl.name, r.score, r.time, r.difficult
		FROM Player as pl
		join Record as r 
		on r.player_id = pl.id
		where r.difficult = 0 and password like $1 and name like $2`

	row := tx.QueryRow(sqlStatement, password, field)

	user = &models.UserPublicInfo{}
	err = row.Scan(&id, &user.Name, &user.BestScore,
		&user.BestTime, &user.Difficult)
	return
}

// GetPrivateInfo get player's personal info
func (db DataBase) getPrivateInfo(tx *sql.Tx, userID int32) (user *models.UserPrivateInfo, err error) {
	sqlStatement := "SELECT name, password " +
		"FROM Player where id = $1"

	row := tx.QueryRow(sqlStatement, userID)

	user = &models.UserPrivateInfo{}
	user.ID = int(userID)
	err = row.Scan(&user.Name, &user.Password)
	if err != nil {
		utils.Debug(true, "cant get user's name and password")
	}

	return
}

// GetUsers returns information about users
// for leaderboard
func (db *DataBase) getUsers(tx *sql.Tx, difficult int, offset int, limit int,
	sort string) (players []*models.UserPublicInfo, err error) {

	sqlStatement := `
	SELECT P.id, P.photo_title, P.name,
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
		return
	}
	defer rows.Close()

	for rows.Next() {
		player := &models.UserPublicInfo{}
		if err = rows.Scan(&player.ID, &player.FileKey, &player.Name, &player.BestScore,
			&player.BestTime, &player.Difficult); err != nil {
			utils.Debug(true, "catch wrong row about one of users")
			break
		}

		players = append(players, player)
	}

	return
}

// GetUser returns information about user
func (db *DataBase) getUser(tx *sql.Tx, userID int32, difficult int) (player *models.UserPublicInfo, err error) {

	sqlStatement := `
	SELECT P.id, P.photo_title, P.name,
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
		&player.BestScore, &player.BestTime, &player.Difficult)
	if err != nil {
		utils.Debug(true, "cant get user")
	}

	return
}

// deletePlayer delete all information about user
func (db *DataBase) deletePlayer(tx *sql.Tx, user *models.UserPrivateInfo) error {
	sqlStatement := `
	DELETE FROM Player where name=$1 and password=$2
	RETURNING ID
		`
	row := tx.QueryRow(sqlStatement, user.Name, user.Password)

	fmt.Println("deletePlayer:", user.Name, user.Password)

	err := row.Scan(&user.ID)
	if err != nil {
		utils.Debug(true, "cant delete player")
	}

	return err
}

// GetPlayerIDbyName get user's id by his name
func (db *DataBase) GetPlayerIDbyName(username string) (id int, err error) {
	sqlStatement := `SELECT id FROM Player WHERE name = $1`
	row := db.Db.QueryRow(sqlStatement, username)

	err = row.Scan(&id)
	if err != nil {
		utils.Debug(false, "cant get player's ID by his name")
	}
	return
}
