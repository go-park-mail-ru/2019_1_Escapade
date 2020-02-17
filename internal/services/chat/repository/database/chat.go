package database

import (
	//
	"context"
	"errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/models"
)

type Chat struct {
	db     infrastructure.Execer
	logger infrastructure.Logger
	trace  infrastructure.ErrorTrace
}

func NewChat(
	db infrastructure.Execer,
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
) (*Chat, error) {
	// check database interface given
	if db == nil {
		return nil, errors.New(ErrNoDatabase)
	}
	// overriding nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}
	// overriding nil value of ErrorTrace
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}
	return &Chat{
		db:     db,
		logger: logger,
		trace:  trace,
	}, nil
}

func (rep *Chat) Get(
	ctx context.Context,
	chatModel *models.Chat,
) (int32, error) {
	if chatModel == nil {
		return 0, rep.trace.New(models.ErrNoChat)
	}
	var id int32
	query := `select id from Chat 
				where chat_type = $1 and type_id = $2;`

	err := rep.db.QueryRowContext(
		ctx,
		query,
		chatModel.Type,
		chatModel.TypeID,
	).Scan(&id)
	return id, err
}

func (rep *Chat) Create(
	ctx context.Context,
	chatType, typeID int32,
) (int32, error) {
	var id int32
	sqlInsert := `INSERT INTO Chat(chat_type, type_id) 
						 VALUES ($1, $2) returning id;`

	err := rep.db.QueryRowContext(
		ctx,
		sqlInsert,
		chatType,
		typeID,
	).Scan(&id)
	return id, err
}
