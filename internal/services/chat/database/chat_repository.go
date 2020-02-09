package database

import (
	//
	_ "github.com/lib/pq"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/infrastructure"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/services/chat/proto"
)

type ChatRepositoryPQ struct{}

func (db *ChatRepositoryPQ) get(tx infrastructure.Transaction,
	ch *proto.Chat) (*proto.ChatID, error) {
	var id int32
	query := `select id from Chat 
					 where chat_type = $1 and type_id = $2;`

	err := tx.QueryRow(query, ch.Type, ch.TypeId).Scan(&id)
	return &proto.ChatID{Value: id}, err
}

func (db *ChatRepositoryPQ) create(tx database.TransactionI,
	chatType, typeID int32) (*proto.ChatID, error) {
	var id int32
	sqlInsert := `INSERT INTO Chat(chat_type, type_id) 
						 VALUES ($1, $2) returning id;`

	err := tx.QueryRow(sqlInsert, chatType, typeID).Scan(&id)
	return &proto.ChatID{Value: id}, err
}
