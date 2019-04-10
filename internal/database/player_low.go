package database

import (
	"database/sql"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"fmt"
	"time"
)

// Check function
type Check func(tx *sql.Tx, value string) error

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

func (db *DataBase) updatePlayerPersonalInfo(tx *sql.Tx, user *models.UserPrivateInfo, oldName string) (id int, err error) {
	sqlStatement := `
			UPDATE Player 
			SET name = $1, email = $2, password = $3, lastSeen = $4
			WHERE name like $5
			RETURNING id
		`

	row := tx.QueryRow(sqlStatement, user.Name,
		user.Email, user.Password, time.Now(), oldName)
	err = row.Scan(&id)
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

func (db *DataBase) checkParameter(tx *sql.Tx, old string, new string, check Check) (choosen string, err error) {
	choosen = old
	if new != old && new != "" {
		if err = check(tx, new); err != nil {
			return
		}
		choosen = new
	}
	return
}

// isNameUnique checks if there are Players with
// this('taken') name and returns corresponding error if yes
func (db DataBase) isNameUnique(tx *sql.Tx, taken string) error {
	sqlStatement := "SELECT name " +
		"FROM Player where name=$1"

	row := tx.QueryRow(sqlStatement, taken)

	var tmp string
	if err := row.Scan(&tmp); err != sql.ErrNoRows {
		if err == nil {
			return re.ErrorNameIstaken()
		}
		return err
	}
	return nil
}

// isEmailUnique checks if there are Players with
// this('taken') email and returns corresponding error if yes
func (db DataBase) isEmailUnique(tx *sql.Tx, taken string) error {
	sqlStatement := "SELECT name " +
		"FROM Player where email=$1"

	row := tx.QueryRow(sqlStatement, taken)

	var tmp string
	if err := row.Scan(&tmp); err != sql.ErrNoRows {
		if err == nil {
			return re.ErrorEmailIstaken()
		}
		return err
	}
	return nil
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

// GetNameByEmail get player's name by his email
func (db DataBase) GetPasswordEmailByName(tx *sql.Tx, name string) (email string, password string, err error) {
	sqlStatement := "SELECT email, password " +
		"FROM Player where name like $1"

	row := tx.QueryRow(sqlStatement, name)

	if err = row.Scan(&email, &password); err != nil {
		return
	}
	return
}

// confirmUnique confirm that user.Email and user.Name
// dont use by another Player
func (db DataBase) confirmUnique(tx *sql.Tx, user *models.UserPrivateInfo) (err error) {
	if err = db.isEmailUnique(tx, user.Email); err != nil {
		return
	}
	err = db.isNameUnique(tx, user.Name)
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
		if err = rows.Scan(&player.ID, &player.FileName, &player.Name, &player.Email, &player.BestScore,
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
