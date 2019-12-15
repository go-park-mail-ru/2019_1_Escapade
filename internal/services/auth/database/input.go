package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/database"
)

type Input struct {
	database.Input
}

func (input *Input) InitAsPSQL() *Input {
	input.Input.InitAsPSQL()
	return input
}
