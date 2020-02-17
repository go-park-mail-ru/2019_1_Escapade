package database

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
)

// User implements the interface UserRepositoryI using database
type User struct {
	db    infrastructure.Execer
	trace infrastructure.ErrorTrace
}

func NewUser(
	dbI infrastructure.Execer,
	trace infrastructure.ErrorTrace,
) (*User, error) {
	if dbI == nil {
		return nil, errors.New(ErrNoDatabase)
	}
	return &User{
		db:    dbI,
		trace: trace,
	}, nil
}

// Create user
func (db *User) Create(
	ctx context.Context,
	user *models.UserPrivateInfo,
) (int, error) {
	if user == nil {
		return 0, db.trace.New(InvalidUser)
	}
	var (
		sqlInsert = `
			INSERT INTO Player(name, password, firstSeen, lastSeen) VALUES
				($1, $2, $3, $4)
				RETURNING id;`
		t  = time.Now()
		id int
	)
	err := db.db.QueryRowContext(
		ctx,
		sqlInsert,
		user.Name,
		user.Password,
		t,
		t,
	).Scan(&id)
	return id, err
}

// Delete delete all information about user
func (db *User) Delete(
	ctx context.Context,
	user *models.UserPrivateInfo,
) error {
	if user == nil {
		return db.trace.New(InvalidUser)
	}
	sqlStatement := `
		DELETE FROM Player where name=$1 and password=$2
			RETURNING ID
		`
	return db.db.QueryRowContext(
		ctx,
		sqlStatement,
		user.Name,
		user.Password,
	).Scan(&user.ID)
}

// UpdateNamePassword update name and password of user with selected id
func (db *User) UpdateNamePassword(
	ctx context.Context,
	user *models.UserPrivateInfo,
) error {
	if user == nil {
		return db.trace.New(InvalidUser)
	}
	sqlStatement := `
			UPDATE Player 
			SET name = $1, password = $2, lastSeen = $3
			WHERE id = $4
			RETURNING id
		`

	return db.db.QueryRowContext(
		ctx,
		sqlStatement,
		user.Name,
		user.Password,
		time.Now(),
		user.ID,
	).Scan(&user.ID)
}

// CheckNamePassword check that there are sych name and password
func (db *User) CheckNamePassword(
	ctx context.Context,
	name, password string,
) (int32, *models.UserPublicInfo, error) {
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
	err := db.db.QueryRowContext(
		ctx,
		sqlStatement,
		password,
		name,
	).Scan(
		&id,
		&user.Name,
		&user.BestScore,
		&user.BestTime,
		&user.Difficult,
	)
	return id, user, err
}

// FetchNamePassword get player's personal info
func (db *User) FetchNamePassword(
	ctx context.Context,
	userID int32,
) (*models.UserPrivateInfo, error) {
	sqlStatement := "SELECT name, password FROM Player where id = $1"

	user := &models.UserPrivateInfo{}
	user.ID = int(userID)
	err := db.db.QueryRowContext(
		ctx,
		sqlStatement,
		userID,
	).Scan(&user.Name, &user.Password)
	return user, err
}

// UpdateLastSeen update users last date seen
func (db *User) UpdateLastSeen(
	ctx context.Context,
	id int,
) error {
	var (
		sqlStatement = `
			UPDATE Player 
			SET lastSeen = $1
			WHERE id = $2
		`
		err error
	)
	_, err = db.db.ExecContext(
		ctx,
		sqlStatement,
		time.Now(),
		id,
	)
	return err
}

// FetchAll returns information about users
func (db *User) FetchAll(
	ctx context.Context,
	params api.UsersSelectParams,
) ([]*models.UserPublicInfo, error) {

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

	return db.fetchAll(ctx, params, sqlStatement)
}

// fetchAll returns information about users
func (db *User) fetchAll(
	ctx context.Context,
	params api.UsersSelectParams,
	sqlStatement string,
) ([]*models.UserPublicInfo, error) {

	players := make([]*models.UserPublicInfo, 0, params.Limit)
	rows, err := db.db.QueryContext(
		ctx,
		sqlStatement,
		params.Difficult,
		params.Offset,
		params.Limit,
	)
	if err != nil {
		return players, err
	}
	defer rows.Close()

	for rows.Next() {
		player := &models.UserPublicInfo{}
		err = rows.Scan(
			&player.ID,
			&player.FileKey,
			&player.Name,
			&player.BestScore,
			&player.BestTime,
			&player.Difficult,
		)
		if err != nil {
			break
		}
		players = append(players, player)
	}

	return players, err
}

// FetchOne returns information about user
func (db *User) FetchOne(
	ctx context.Context,
	userID int32,
	difficult int,
) (*models.UserPublicInfo, error) {

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
	err := db.db.QueryRowContext(ctx,
		sqlStatement,
		userID,
		difficult,
	).Scan(
		&player.ID,
		&player.FileKey,
		&player.Name,
		&player.BestScore,
		&player.BestTime,
		&player.Difficult,
	)

	return player, err
}

// PagesCount return user's pages count
func (db *User) PagesCount(
	ctx context.Context,
	perPage int,
) (int, error) {
	sqlStatement := `SELECT count(1) FROM Player`
	var amount int
	err := db.db.QueryRowContext(
		ctx,
		sqlStatement,
	).Scan(&amount)
	if err != nil {
		return 0, err
	}
	pageUsers := 10 // в конфиг
	amount = db.fixAmount(amount, pageUsers)
	perPage = db.fixPerPage(perPage)
	amount = int(math.Ceil(float64(amount) / float64(perPage)))
	return amount, nil
}

func (db *User) fixAmount(amount, pageUsers int) int {
	if amount > pageUsers {
		return pageUsers
	}
	return amount
}

func (db *User) fixPerPage(perPage int) int {
	if perPage <= 0 {
		return 1
	}
	return perPage
}
