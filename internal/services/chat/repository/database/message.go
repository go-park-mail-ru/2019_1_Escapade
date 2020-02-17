package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/models/database"
)

type Message struct {
	db     infrastructure.Execer
	logger infrastructure.Logger
	trace  infrastructure.ErrorTrace
	photo  infrastructure.PhotoService
}

func NewMessage(
	db infrastructure.Execer,
	photo infrastructure.PhotoService,
	logger infrastructure.Logger,
	trace infrastructure.ErrorTrace,
) (*Message, error) {
	// check database interface given
	if db == nil {
		return nil, errors.New(ErrNoDatabase)
	}
	// overriding nil value of PhotoService
	if photo == nil {
		photo = new(infrastructure.PhotoServiceNil)
	}
	// overriding nil value of Logger
	if logger == nil {
		logger = new(infrastructure.LoggerNil)
	}
	// overriding nil value of ErrorTrace
	if trace == nil {
		trace = new(infrastructure.ErrorTraceNil)
	}
	return &Message{
		db:     db,
		logger: logger,
		trace:  trace,
		photo:  photo,
	}, nil
}

func (rep *Message) CreateOne(
	ctx context.Context,
	msg *models.Message,
) (int32, error) {
	if msg == nil {
		return 0, rep.trace.New(models.ErrNoMessage)
	}
	var (
		extraParameters int
		query           = "INSERT INTO Message(not_saved_id"
	)
	if msg.Answer != nil {
		query += ",answer_id"
		extraParameters++
	}
	if msg.To != nil {
		query += ",getter_id, getter_name"
		extraParameters += 2
	}

	query += `,sender_id, sender_name, sender_status, 
		chat_id, message, date) VALUES
		($1,$2,$3,$4,$5,$6,$7`
	var row *sql.Row
	switch extraParameters {
	case 0:
		row = rep.createMsgNoAnswerOrGetter(ctx, query, msg)
	case 1:
		row = rep.createMsgWithAnswer(ctx, query, msg)
	case 2:
		row = rep.createMsgWithGetter(ctx, query, msg)
	case 3:
		row = rep.createMsgWithGetterAndAnswer(ctx, query, msg)
	}
	if row == nil {
		return 0, rep.trace.New(ErrInternal)
	}
	var id int32
	err := row.Scan(&id)
	return id, err
}

func (rep *Message) createMsgNoAnswerOrGetter(
	ctx context.Context,
	query string,
	message *models.Message,
) *sql.Row {
	return rep.db.QueryRowContext(
		ctx,
		query+`) returning id;`,
		message.ID,
		message.From.ID,
		message.From.Name,
		message.From.Status,
		message.ChatID,
		message.Text,
		message.Time,
	)
}

func (rep *Message) createMsgWithAnswer(
	ctx context.Context,
	query string,
	message *models.Message,
) *sql.Row {
	return rep.db.QueryRowContext(
		ctx,
		query+`$8) returning id;`,
		message.ID,
		message.Answer.ID,
		message.From.ID,
		message.From.Name,
		message.From.Status,
		message.ChatID,
		message.Text,
		message.Time,
	)
}

func (rep *Message) createMsgWithGetter(
	ctx context.Context,
	query string,
	message *models.Message,
) *sql.Row {
	return rep.db.QueryRowContext(
		ctx,
		query+`$8, $9) returning id;`,
		message.ID,
		message.To.ID,
		message.To.Name,
		message.From.ID,
		message.From.Name,
		message.From.Status,
		message.ChatID,
		message.Text,
		message.Time,
	)
}

func (rep *Message) createMsgWithGetterAndAnswer(
	ctx context.Context,
	query string,
	message *models.Message,
) *sql.Row {
	return rep.db.QueryRowContext(
		ctx,
		query+`$8, $9, $10) returning id;`,
		message.ID,
		message.Answer.ID,
		message.To.ID,
		message.To.Name,
		message.From.ID,
		message.From.Name,
		message.From.Status,
		message.ChatID,
		message.Text,
		message.Time,
	)
}

func (rep *Message) CreateMany(
	ctx context.Context,
	messages *models.Messages,
) ([]int32, error) {
	if len(messages.Messages) == 0 {
		return nil, nil
	}
	if messages == nil {
		return nil, rep.trace.New(models.ErrNoMessages)
	}

	var (
		err error
		ids = make([]int32, len(messages.Messages))
	)
	for i, message := range messages.Messages {
		ids[i], err = rep.CreateOne(ctx, message)
		if err != nil {
			break
		}
	}
	return ids, err
}

func (rep *Message) Update(
	ctx context.Context,
	message *models.Message,
) error {
	var (
		id        int32
		sqlUpdate = `
	Update Message set message = $1, edited = true where 
		id = $2 or (not_saved_id = $2 and sender_id = $3)
		RETURNING ID;`

		err = rep.db.QueryRowContext(
			ctx,
			sqlUpdate,
			message.Text,
			message.ID,
			message.From.ID,
		).Scan(&id)
	)

	return err
}

func (rep *Message) Delete(
	ctx context.Context,
	message *models.Message,
) error {
	var (
		id        int32
		sqlDelete = `
		Delete from Message 
			where id = $1 or 
				(not_saved_id = $1 and sender_id = $2)
			RETURNING ID;`

		err = rep.db.QueryRowContext(
			ctx,
			sqlDelete,
			message.ID,
			message.From.ID,
		).Scan(&id)
	)
	return err
}

func (rep *Message) GetAll(
	ctx context.Context,
	chatID int32,
) ([]*models.Message, error) {

	sqlStatement := `
	select M.id, M.answer_id, M.sender_id, M.sender_name, M.sender_status, 
		M.getter_id, M.getter_name, M.chat_id, M.message, M.date, M.edited,
		S.photo_title, A.sender_id, A.sender_name, A.sender_status,
		A.message, A.getter_id, A.getter_name, A.chat_id, A.date
		from Message as M 
		left join Player as S on M.sender_id = S.id
		left join Message as A on M.answer_id = A.id
		where M.chat_id = $1
		ORDER BY M.ID ASC;
		`
	rows, err := rep.db.QueryContext(
		ctx,
		sqlStatement,
		chatID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	messages := make([]*models.Message, 0)

	for rows.Next() {
		var (
			answer  = database.Message{}
			message = database.Message{}
			photo   sql.NullString
		)
		err = rows.Scan(
			&message.ID,
			&answer.ID,
			&message.From.ID,
			&message.From.Name,
			&message.From.Status,
			&message.To.ID,
			&message.To.Name,
			&message.ChatID,
			&message.Text,
			&message.Time,
			&message.Edited,
			&photo,
			&answer.From.ID,
			&answer.From.Name,
			&answer.From.Status,
			&answer.Text,
			&answer.To.ID,
			&answer.To.Name,
			&answer.ChatID,
			(*database.NullTime)(&answer.Time))
		if err != nil {
			break
		}
		result := message.Get()
		if photo.Valid {
			result.From.Photo = photo.String
		} else {
			result.From.Photo = rep.photo.GetDefaultAvatar()
		}
		rep.logger.Println(
			"result.From.Photo",
			result.From.Photo,
		)

		messages = append(messages, result)
	}
	if err != nil {
		rep.logger.Println("cant get message:", err.Error())
		return nil, err
	}

	return messages, nil
}

// 370 -> 319
