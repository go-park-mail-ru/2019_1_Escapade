package chat

import (
	"database/sql"

	//
	_ "github.com/lib/pq"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

func (service *Service) insertMessage(message *Message) (*MessageID, error) {

	var (
		id  int32
		err error
	)
	sqlInsert := `
	INSERT INTO Message(answer_id, sender_id, sender_name, sender_status, 
		getter_id, getter_name,chat_id, message, date) VALUES
		`
	sqlInsert += addMessageToQuery(message) + " returning id;"
	row := service.DB.QueryRow(sqlInsert)

	if err = row.Scan(&id); err != nil {
		utils.Debug(true, "cant create message", err.Error())
		return &MessageID{}, err
	}

	return &MessageID{Value: id}, nil
}

func (service *Service) insertMessages(messages *Messages) (*MessagesID, error) {

	sqlInsert := `
	INSERT INTO Message(answer_id, sender_id, sender_name, sender_status,
		 getter_id, getter_name, chat_id, message, date) VALUES
		`

	if len(messages.Messages) == 0 {
		return &MessagesID{}, nil
	}
	for i, message := range messages.Messages {
		if i == 0 {
			sqlInsert += addMessageToQuery(message)
		} else {
			sqlInsert += "," + addMessageToQuery(message)
		}
	}

	sqlInsert += " returning id;"
	var (
		rows *sql.Rows
		err  error
	)

	if rows, err = service.DB.Query(sqlInsert); err != nil {
		return &MessagesID{}, err
	}
	defer rows.Close()

	var ids = make([]int32, len(messages.Messages))
	i := 0
	for rows.Next() {
		if err = rows.Scan(&ids[i]); err != nil {
			break
		}
		i++
	}

	if err != nil {
		return &MessagesID{}, err
	}

	return &MessagesID{Values: ids}, nil
}

func (service *Service) updateMessage(message *Message) (*Result, error) {

	var (
		id  int32
		err error
	)
	sqlUpdate := `
	Update Message set message = $1, edited = true where id = $2
		RETURNING ID;
		`
	row := service.DB.QueryRow(sqlUpdate, message.Text, message.Id)

	if err = row.Scan(&id); err != nil {
		utils.Debug(true, "cant update message", err.Error())
		return &Result{Done: false}, err
	}

	return &Result{Done: true}, nil
}

func (service *Service) deleteMessage(message *Message) (*Result, error) {

	var (
		id  int32
		err error
	)
	sqlDelete := `
	Delete from Message where id = $1
	RETURNING ID;
		`
	row := service.DB.QueryRow(sqlDelete, message.Id)

	if err = row.Scan(&id); err != nil {
		utils.Debug(true, "cant delete message", err.Error())
		return &Result{Done: false}, err
	}

	return &Result{Done: true}, nil
}

func (service *Service) getChatMessages(chatID *ChatID) (*Messages, error) {

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
	rows, err = service.DB.Query(sqlStatement, chatID.Value)
	if err != nil {
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
			(*models.ScanTime)(&answer.Time)); err != nil {
			break
		}
		result = MessageFromNullMessage(message)
		result.Answer = MessageFromNullMessage(answer)
		result.From.Photo = "anonymous.jpg"

		messages = append(messages, result)
	}
	if err != nil {
		utils.Debug(true, "cant get message:", err.Error())
		return &Messages{}, err
	}

	return &Messages{Messages: messages}, err
}
