package database

import (
	//
	_ "github.com/lib/pq"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/utils"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

type UserRepositoryPQ struct{}

func addUserToQuery(user *proto.User) string {
	return "('" + utils.String32(user.Id) + "',$1)"
}

func (db *UserRepositoryPQ) create(tx database.TransactionI, chatID int32,
	users ...*proto.User) error {
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

func (db *UserRepositoryPQ) delete(Db database.DatabaseI,
	userInGroup *proto.UserInGroup) (*proto.Result, error) {

	var (
		id  int32
		err error
	)
	sqlDelete := `
	Delete from UserInChat where user_id = $1 and chat_id = $2;
		`
	row := Db.QueryRow(sqlDelete, userInGroup.User.Id, userInGroup.Chat.Id)

	if err = row.Scan(&id); err != nil {
		utils.Debug(true, "cant delete message", err.Error())
		return &proto.Result{Done: false}, err
	}

	return &proto.Result{Done: true}, nil
}
