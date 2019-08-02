package chat

import (
	"database/sql"

	//
	_ "github.com/lib/pq"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

func (service *Service) insertUsers(tx *sql.Tx, chatID int32, users ...*User) error {
	var (
		err error
	)
	sqlInsert := `INSERT INTO UserInChat(user_id, chat_id) VALUES `

	if len(users) == 0 {
		return nil
	}
	for i, user := range users {
		if i == 0 {
			sqlInsert += addUserToQuery(user)
		} else {
			sqlInsert += "," + addUserToQuery(user)
		}
	}

	_, err = tx.Exec(sqlInsert, chatID)

	return err
}

func (service *Service) deleteUserInChat(userInGroup *UserInGroup) (*Result, error) {

	var (
		id  int32
		err error
	)
	sqlDelete := `
	Delete from UserInChat where user_id = $1 and chat_id = $2;
		`
	row := service.DB.QueryRow(sqlDelete, userInGroup.User.Id, userInGroup.Chat.Id)

	if err = row.Scan(&id); err != nil {
		utils.Debug(true, "cant delete message", err.Error())
		return &Result{Done: false}, err
	}

	return &Result{Done: true}, nil
}
