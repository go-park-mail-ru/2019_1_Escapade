package database

import (
	"database/sql"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"fmt"
)

// createPlayer create player
func (db *DataBase) createMessage(tx *sql.Tx, message *models.Message,
	inRoom bool, gameID string) (err error) {
	sqlInsert := `
	INSERT INTO GameChat(player_id, in_room, roomID, message, time) VALUES
		($1, $2, $3, $4, $5);
		`
	_, err = db.Db.Exec(sqlInsert, message.User.ID, inRoom,
		gameID, message.Text, message.Time)

	if err != nil {
		fmt.Println("createMessage err:", err.Error())
		return
	}
	fmt.Println("createMessage success")

	return
}

// GetUsers returns information about users
// for leaderboard
func (db *DataBase) getMessages(tx *sql.Tx, inRoom bool, gameID string) (messages []*models.Message, err error) {

	var (
		rows *sql.Rows
	)
	sqlStatement := `
	select GC.player_id, P.name, P.photo_title, GC.message, GC.time 
		from GameChat as GC 
		join Player as P on P.id = GC.player_id`
	if inRoom {
		sqlStatement += ` where GC.roomID like $1;`
		rows, err = tx.Query(sqlStatement, gameID)
	} else {
		sqlStatement += ` where GC.in_room = false;`
		rows, err = tx.Query(sqlStatement)
	}
	if err != nil {
		fmt.Println("database/getMessages cant access to database:", err.Error())
		return
	}

	defer rows.Close()
	messages = make([]*models.Message, 0)

	for rows.Next() {
		user := &models.UserPublicInfo{}
		message := &models.Message{
			User: user,
		}

		if err = rows.Scan(&message.User.ID, &message.User.Name,
			&message.User.PhotoURL, &message.Text, &message.Time); err != nil {

			break
		}
		fmt.Println("load message:", message)

		messages = append(messages, message)
	}
	if err != nil {
		fmt.Println("database/GetUsers wrong row catched:", err.Error())
		return
	}

	fmt.Println("database/getMessages +")

	return
}
