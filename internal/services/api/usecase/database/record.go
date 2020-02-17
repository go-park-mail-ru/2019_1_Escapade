package database

import (
	"context"
	"errors"
	"time"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/api/repository/database"
)

// Record implements the interface RecordUseCaseI
type Record struct {
	db             infrastructure.Database
	trace          infrastructure.ErrorTrace
	recordDB       api.RecordRepositoryI
	contextTimeout time.Duration
}

// NewRecord create new instance of Record
func NewRecord(
	dbI infrastructure.Database,
	trace infrastructure.ErrorTrace,
	timeout time.Duration,
) (*Record, error) {
	if dbI == nil {
		return nil, errors.New(ErrNoDatabase)
	}
	recordRep, err := database.NewRecord(dbI, trace)
	if err != nil {
		return nil, err
	}
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}
	return &Record{
		db:             dbI,
		trace:          trace,
		recordDB:       recordRep,
		contextTimeout: timeout,
	}, nil
}

// Update records for offline game
func (usecase *Record) Update(
	c context.Context,
	id int32,
	record *models.Record,
) error {
	if record == nil {
		return usecase.trace.New(InvalidRecord)
	}
	if id <= 0 {
		return usecase.trace.New(InvalidID)
	}
	ctx, cancel := context.WithTimeout(
		c,
		usecase.contextTimeout,
	)
	defer cancel()
	err := usecase.recordDB.Update(ctx, id, record)
	return err
}
