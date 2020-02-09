package database

import (
	"context"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
)

// Image implements the interface ImageRepositoryI using the sql postgres driver
type Image struct {
	db infrastructure.Execer
}

func NewImage(dbI infrastructure.Execer) *Image {
	return &Image{dbI}
}

// Update set filename of avatar to relation Player
func (image *Image) Update(
	ctx context.Context,
	filename string,
	userID int32,
) error {
	sqlStatement := `UPDATE Player SET photo_title = $1 WHERE id = $2;`
	_, err := image.db.ExecContext(
		ctx,
		sqlStatement,
		filename,
		userID,
	)
	return err
}

// FetchByName Get avatar - filename of player image by his name
func (image *Image) FetchByName(
	ctx context.Context,
	name string,
) (string, error) {
	sqlStatement := `SELECT photo_title FROM Player WHERE name like $1`
	var filename string
	err := image.db.QueryRowContext(
		ctx,
		sqlStatement,
		name,
	).Scan(&filename)
	return filename, err
}

// FetchByID Get avatar - filename of player image by his id
func (image *Image) FetchByID(
	ctx context.Context,
	id int32,
) (string, error) {
	sqlStatement := `SELECT photo_title FROM Player WHERE id=$1`
	var filename string
	err := image.db.QueryRowContext(
		ctx,
		sqlStatement,
		id,
	).Scan(&filename)
	return filename, err
}

// 33
