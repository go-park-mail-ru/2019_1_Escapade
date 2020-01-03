package database

import (
	"math"
	"time"

	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
)

// UserRepositoryPQ implements the interface UserRepositoryI using the sql postgres driver
type UserRepositoryPQ struct{}

// UsersSelectParams parameters to select user
type UsersSelectParams struct {
	Difficult int
	Offset    int
	Limit     int
	Sort      string
}

// Create user
func (db *UserRepositoryPQ) Create(tx idb.TransactionI, user *models.UserPrivateInfo) (int, error) {
	var (
		sqlInsert = `
			INSERT INTO Player(name, password, firstSeen, lastSeen) VALUES
				($1, $2, $3, $4)
				RETURNING id;`
		t  = time.Now()
		id int
	)
	err := tx.QueryRow(sqlInsert, user.Name, user.Password, t, t).Scan(&id)
	return id, err
}

// Delete delete all information about user
func (db *UserRepositoryPQ) Delete(tx idb.TransactionI, user *models.UserPrivateInfo) error {
	sqlStatement := `
	DELETE FROM Player where name=$1 and password=$2
	RETURNING ID
		`
	return tx.QueryRow(sqlStatement, user.Name, user.Password).Scan(&user.ID)
}

// UpdateNamePassword update name and password of user with selected id
func (db *UserRepositoryPQ) UpdateNamePassword(tx idb.TransactionI, user *models.UserPrivateInfo) error {
	sqlStatement := `
			UPDATE Player 
			SET name = $1, password = $2, lastSeen = $3
			WHERE id = $4
			RETURNING id
		`

	return tx.QueryRow(sqlStatement, user.Name, user.Password, time.Now(), user.ID).Scan(&user.ID)
}

// CheckNamePassword check that there are sych name and password
func (db *UserRepositoryPQ) CheckNamePassword(tx idb.TransactionI, name string, password string) (int32, *models.UserPublicInfo, error) {
	var (
		sqlStatement = `
			SELECT pl.id, pl.name, r.score, r.time, r.difficult
			FROM Player as pl
			join Record as r 
			on r.player_id = pl.id
			where r.difficult = 0 and password like $1 and name like $2`
		id int32
	)
	user := &models.UserPublicInfo{}
	err := tx.QueryRow(sqlStatement, password, name).Scan(&id, &user.Name, &user.BestScore, &user.BestTime, &user.Difficult)
	return id, user, err
}

// FetchNamePassword get player's personal info
func (db *UserRepositoryPQ) FetchNamePassword(tx idb.TransactionI, userID int32) (*models.UserPrivateInfo, error) {
	sqlStatement := "SELECT name, password FROM Player where id = $1"

	user := &models.UserPrivateInfo{}
	user.ID = int(userID)
	err := tx.QueryRow(sqlStatement, userID).Scan(&user.Name, &user.Password)

	return user, err
}

// UpdateLastSeen update users last date seen
func (db *UserRepositoryPQ) UpdateLastSeen(tx idb.TransactionI, id int) error {
	var (
		sqlStatement = `
			UPDATE Player 
			SET lastSeen = $1
			WHERE id = $2
		`
		err error
	)
	_, err = tx.Exec(sqlStatement, time.Now(), id)
	return err
}

// FetchAll returns information about users
func (db *UserRepositoryPQ) FetchAll(tx idb.TransactionI, params UsersSelectParams) ([]*models.UserPublicInfo, error) {

	sqlStatement := `
	SELECT P.id, P.photo_title, P.name,
				 R.score, R.time, R.Difficult
	FROM Player as P
	join Record as R 
	on R.player_id = P.id
	where r.difficult = $1  
	`
	if params.Sort == "score" {
		sqlStatement += ` ORDER BY (score) desc `
	} else {
		sqlStatement += ` ORDER BY (time) `
	}
	sqlStatement += ` OFFSET $2 Limit $3 `

	return db.fetchAll(tx, params, sqlStatement)
}

// fetchAll returns information about users
func (db *UserRepositoryPQ) fetchAll(tx idb.TransactionI, params UsersSelectParams,
	sqlStatement string) ([]*models.UserPublicInfo, error) {

	players := make([]*models.UserPublicInfo, 0, params.Limit)
	rows, err := tx.Query(sqlStatement, params.Difficult, params.Offset,
		params.Limit)
	if err != nil {
		return players, err
	}
	defer rows.Close()

	for rows.Next() {
		player := &models.UserPublicInfo{}
		err = rows.Scan(&player.ID, &player.FileKey, &player.Name,
			&player.BestScore, &player.BestTime, &player.Difficult)
		if err != nil {
			break
		}
		players = append(players, player)
	}

	return players, err
}

// FetchOne returns information about user
func (db *UserRepositoryPQ) FetchOne(tx idb.TransactionI, userID int32,
	difficult int) (*models.UserPublicInfo, error) {

	sqlStatement := `
	SELECT P.id, P.photo_title, P.name,
				 R.score, R.time, R.Difficult
	FROM Player as P
	join Record as R 
	on R.player_id = P.id
	where R.player_id = $1 and
		R.difficult = $2
	`

	player := &models.UserPublicInfo{}
	err := tx.QueryRow(sqlStatement, userID, difficult).Scan(
		&player.ID, &player.FileKey, &player.Name,
		&player.BestScore, &player.BestTime, &player.Difficult)

	return player, err
}

// PagesCount return user's pages count
func (db *UserRepositoryPQ) PagesCount(dbI idb.Interface, perPage int) (int, error) {
	sqlStatement := `SELECT count(1) FROM Player`
	var amount int
	err := dbI.QueryRow(sqlStatement).Scan(&amount)
	if err != nil {
		return 0, err
	}
	pageUsers := 10 // в конфиг
	amount = db.fixAmount(amount, pageUsers)
	perPage = db.fixPerPage(perPage)
	amount = int(math.Ceil(float64(amount) / float64(perPage)))
	return amount, nil
}

func (db *UserRepositoryPQ) fixAmount(amount, pageUsers int) int {
	if amount > pageUsers {
		return pageUsers
	}
	return amount
}

func (db *UserRepositoryPQ) fixPerPage(perPage int) int {
	if perPage <= 0 {
		return 1
	}
	return perPage
}
