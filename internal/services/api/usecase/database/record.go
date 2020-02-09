package database

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/repository/database"
)

// Record implements the interface RecordUseCaseI
type Record struct {
	db             infrastructure.Database
	recordDB       api.RecordRepositoryI
	contextTimeout time.Duration
}

func NewRecord(
	dbI infrastructure.Database,
	timeout time.Duration,
) *Record {
	return &Record{
		db:             dbI,
		recordDB:       database.NewRecord(dbI),
		contextTimeout: timeout,
	}
}

// Update records for offline game
func (usecase *Record) Update(
	c context.Context,
	id int32,
	record *models.Record,
) error {
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()
	err := usecase.recordDB.Update(ctx, id, record)
	return err
}
