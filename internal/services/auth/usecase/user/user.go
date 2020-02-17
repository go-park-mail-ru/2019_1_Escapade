package user

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// User implements the interface UserRepositoryI using database
type UserDB struct {
	db infrastructure.Execer
}

func NewUserDB(dbI infrastructure.Execer) *UserDB {
	return &UserDB{dbI}
}

// CheckNamePassword check that there are sych name and password
func (db *UserDB) CheckNamePassword(
	ctx context.Context,
	name, password string,
) (int32, error) {
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
	return id, err
}
