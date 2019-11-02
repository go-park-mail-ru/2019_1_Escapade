package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
)

// RecordUseCase implements the interface RecordUseCaseI
type RecordUseCase struct {
	UseCaseBase
	record RecordRepositoryI
}

func (db *RecordUseCase) Init(record RecordRepositoryI) {
	db.record = record
}

// UpdateRecords update records for offline game
func (db *RecordUseCase) Update(id int32, record *models.Record) error {
	tx, err := db.Db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = db.record.update(tx, id, record); err != nil {
		return err
	}

	err = tx.Commit()
	return err
}
