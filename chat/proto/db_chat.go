package chat

import (
	"database/sql"

	//
	_ "github.com/lib/pq"

	"github.com/golang/protobuf/ptypes"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

func addMessageToQuery(message *Message) string {
	if message.Answer == nil {
		message.Answer = &Message{}
	}
	if message.From == nil {
		message.From = &User{}
	}
	if message.To == nil {
		message.To = &User{}
	}

	return "('" + utils.String32(message.Answer.Id) + "','" +
		utils.String32(message.From.Id) + "','" + message.From.Name + "','" +
		utils.String32(int32(message.From.Status)) + "','" +
		utils.String32(message.To.Id) + "','" + message.To.Name + "','" +
		utils.String32(message.ChatId) + "','" + message.Text + "','" +
		ptypes.TimestampString(message.Time) + "')"
}

func addUserToQuery(user *User) string {
	return "('" + utils.String32(user.Id) + "',$1)"
}

func (service *Service) getChat(chat *Chat) (*ChatID, error) {

	var (
		id  int32
		err error
	)

	query := `select id from Chat where chat_type = $1 and type_id = $2;`

	row := service.DB.QueryRow(query, chat.Type, chat.TypeId)

	if err = row.Scan(&id); err != nil {
		utils.Debug(true, "cant get chat", err.Error())
		return &ChatID{}, err
	}

	return &ChatID{Value: id}, nil
}

func (service *Service) insertChat(tx *sql.Tx, chat *ChatWithUsers) (*ChatID, error) {
	var (
		id  int32
		err error
	)
	sqlInsert := `
	INSERT INTO Chat(chat_type, type_id) VALUES ($1, $2) returning id;
		`
	row := tx.QueryRow(sqlInsert, chat.Type, chat.TypeId)

	if err = row.Scan(&id); err != nil {
		utils.Debug(true, "cant create message", err.Error())
		return &ChatID{}, err
	}

	return &ChatID{Value: id}, nil
}
