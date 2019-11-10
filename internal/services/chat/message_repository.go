package chat

import (
	"database/sql"
	fmt "fmt"
	"time"

	"github.com/golang/protobuf/ptypes"

	//
	_ "github.com/lib/pq"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

type MessageRepositoryPQ struct{}

func (db *MessageRepositoryPQ) createOne(Db database.DatabaseI, message *Message) (*MessageID, error) {

	var (
		id              int32
		err             error
		extraParameters int
		row             *sql.Row
		date            time.Time
	)
	var sqlInsert = "INSERT INTO Message(not_saved_id"
	if message.Answer != nil {
		sqlInsert += ",answer_id"
		extraParameters++
	}
	if message.To != nil {
		sqlInsert += ",getter_id, getter_name"
		extraParameters += 2
	}

	date, err = ptypes.Timestamp(message.Time)
	if err != nil {
		utils.Debug(false, "cant convert ptypes.Timestamp to time.Time", err.Error())
		return nil, err
	}

	sqlInsert += `,sender_id, sender_name, sender_status, 
	chat_id, message, date) VALUES
	($1,$2,$3,$4,$5,$6,$7`
	switch extraParameters {
	case 0:
		sqlInsert += `) returning id;`
		row = Db.QueryRow(sqlInsert, message.Id, message.From.Id,
			message.From.Name, message.From.Status, message.ChatId,
			message.Text, date)
	case 1:
		sqlInsert += `$8) returning id;`
		row = Db.QueryRow(sqlInsert, message.Id, message.Answer.Id,
			message.From.Id, message.From.Name, message.From.Status,
			message.ChatId, message.Text, date)
	case 2:
		row = Db.QueryRow(sqlInsert, message.Id, message.To.Id,
			message.To.Name, message.From.Id, message.From.Name,
			message.From.Status, message.ChatId, message.Text, date)
		sqlInsert += `$8, $9) returning id;`
	case 3:
		row = Db.QueryRow(sqlInsert, message.Id, message.Answer.Id,
			message.To.Id, message.To.Name, message.From.Id, message.From.Name,
			message.From.Status, message.ChatId, message.Text, date)
		sqlInsert += `$8, $9, $10) returning id;`
	}

	if err = row.Scan(&id); err != nil {
		utils.Debug(false, "sql statement:", sqlInsert)
		utils.Debug(false, "cant create message", err.Error())
		return &MessageID{}, err
	}

	return &MessageID{Value: id}, nil
}

func (db *MessageRepositoryPQ) createMany(Db database.DatabaseI, messages *Messages) (*MessagesID, error) {

	var err error
	if len(messages.Messages) == 0 {
		return &MessagesID{}, nil
	}
	var ids = make([]*MessageID, len(messages.Messages))
	for i, message := range messages.Messages {
		ids[i], err = db.createOne(Db, message)
		if err != nil {
			break
		}
	}

	return &MessagesID{Values: ids}, err
}

func (db *MessageRepositoryPQ) update(Db database.DatabaseI, message *Message) (*Result, error) {

	var (
		id  int32
		err error
	)
	sqlUpdate := `
	Update Message set message = $1, edited = true where id = $2 or (not_saved_id = $2 and sender_id = $3)
		RETURNING ID;
		`
	row := Db.QueryRow(sqlUpdate, message.Text, message.Id, message.From.Id)

	if err = row.Scan(&id); err != nil {
		utils.Debug(false, "sql statement:", sqlUpdate)
		utils.Debug(false, "cant update message", err.Error())
		return &Result{Done: false}, err
	}

	return &Result{Done: true}, nil
}

func (db *MessageRepositoryPQ) delete(Db database.DatabaseI, message *Message) (*Result, error) {

	var (
		id  int32
		err error
	)
	sqlDelete := `
	Delete from Message where id = $1 or (not_saved_id = $1 and sender_id = $2)
	RETURNING ID;
		`
	row := Db.QueryRow(sqlDelete, message.Id, message.From.Id)

	if err = row.Scan(&id); err != nil {
		utils.Debug(false, "sql statement:", sqlDelete)
		utils.Debug(false, "cant delete message", err.Error())
		return &Result{Done: false}, err
	}

	return &Result{Done: true}, nil
}

func (db *MessageRepositoryPQ) getAll(Db database.DatabaseI, chatID *ChatID) (*Messages, error) {

	var (
		rows  *sql.Rows
		err   error
		photo sql.NullString
	)
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
	rows, err = Db.Query(sqlStatement, chatID.Value)
	if err != nil {
		fmt.Println("wtf man", err.Error())
		return &Messages{}, err
	}

	defer rows.Close()
	messages := make([]*Message, 0)

	for rows.Next() {
		var (
			aFrom  = &models.MessageUserSQL{}
			aTo    = &models.MessageUserSQL{}
			answer = &models.MessageSQL{
				From: aFrom,
				To:   aTo,
			}
			from    = &models.MessageUserSQL{}
			to      = &models.MessageUserSQL{}
			message = &models.MessageSQL{
				From: from,
				To:   to,
			}
			result = &Message{}
		)

		if err = rows.Scan(&message.ID, &answer.ID, &from.ID, &from.Name,
			&from.Status, &to.ID, &to.Name, &message.ChatID, &message.Text,
			&message.Time, &message.Edited, &photo, &aFrom.ID, &aFrom.Name,
			&aFrom.Status, &answer.Text, &aTo.ID, &aTo.Name, &answer.ChatID,
			(*models.NullTime)(&answer.Time)); err != nil {
			break
		}
		result, err = MessageFromNullMessage(message)
		if err != nil {
			break
		}
		result.Answer, err = MessageFromNullMessage(answer)
		if err != nil {
			break
		}
		if photo.Valid {
			result.From.Photo = photo.String
		} else {
			result.From.Photo = "anonymous.jpg"
		}
		utils.Debug(false, "result.From.Photo", result.From.Photo)

		messages = append(messages, result)
	}
	if err != nil {
		utils.Debug(true, "cant get message:", err.Error())
		return &Messages{}, err
	}

	return &Messages{Messages: messages}, err
}
