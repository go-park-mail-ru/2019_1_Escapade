package chat

import (

	//

	_ "github.com/lib/pq"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

type ChatRepositoryPQ struct{}

func (db *ChatRepositoryPQ) get(tx database.TransactionI, ch *Chat) (*ChatID, error) {

	var (
		id  int32
		err error
	)

	query := `select id from Chat where chat_type = $1 and type_id = $2;`

	row := tx.QueryRow(query, ch.Type, ch.TypeId)

	if err = row.Scan(&id); err != nil {
		utils.Debug(false, "cant get chaaat", err.Error())
		return &ChatID{}, err
	}

	return &ChatID{Value: id}, nil
}

func (db *ChatRepositoryPQ) create(tx database.TransactionI, chatType, typeID int32) (*ChatID, error) {
	var (
		id  int32
		err error
	)
	sqlInsert := `
	INSERT INTO Chat(chat_type, type_id) VALUES ($1, $2) returning id;
		`
	row := tx.QueryRow(sqlInsert, chatType, typeID)

	if err = row.Scan(&id); err != nil {
		utils.Debug(true, "cant create message", err.Error())
		return &ChatID{}, err
	}

	return &ChatID{Value: id}, nil
}
