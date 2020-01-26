package database

import (
	"context"
	"time"

	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/repository/database"
)

// Record implements the interface RecordUseCaseI
type Record struct {
	db             idb.Interface
	recordDB       api.RecordRepositoryI
	contextTimeout time.Duration
}

func NewRecord(dbI idb.Interface, timeout time.Duration) *Record {
	return &Record{
		db:             dbI,
		recordDB:       database.NewRecord(dbI),
		contextTimeout: timeout,
	}
}

// Update records for offline game
func (repository *Record) Update(c context.Context, id int32, record *models.Record) error {
	ctx, cancel := context.WithTimeout(c, repository.contextTimeout)
	defer cancel()
	err := repository.recordDB.Update(ctx, id, record)
	return err
}
