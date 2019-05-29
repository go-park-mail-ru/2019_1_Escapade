package database

import (
	"database/sql"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"fmt"
)

// GetUsers returns information about users
// for leaderboard
func (db *DataBase) getMessages(tx *sql.Tx, inRoom bool, gameID string) (messages []*models.Message, err error) {

	var (
		rows *sql.Rows
	)
	sqlStatement := `
	select GC.player_id, GC.name, P.name, P.photo_title, GC.message, GC.time 
		from GameChat as GC 
		left join Player as P on P.id = GC.player_id`
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
		userSQL := &models.UserPublicInfoSQL{}
		message := &models.Message{
			User: user,
		}

		if err = rows.Scan(&userSQL.ID, &user.Name, &userSQL.Name,
			&userSQL.PhotoURL, &message.Text, &message.Time); err != nil {

			break
		}
		if id, erro := userSQL.ID.Value(); erro == nil {
			user.ID = int(id.(int64))
		}
		if name, _ := userSQL.Name.Value(); name != nil {
			user.Name = name.(string)
		}
		if photoURL, _ := userSQL.PhotoURL.Value(); photoURL != nil {
			user.PhotoURL = photoURL.(string)
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
