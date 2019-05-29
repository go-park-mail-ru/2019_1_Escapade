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
		left join Player as P on P.id = GC.player_id
		`
	if inRoom {
		sqlStatement += ` where GC.roomID like $1 ORDER BY GC.ID ASC;`
		rows, err = tx.Query(sqlStatement, gameID)
	} else {
		sqlStatement += ` where GC.in_room = false ORDER BY GC.ID ASC;`
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
		user.PhotoURL = "https://escapade.hb.bizmrg.com/76813f9e-60a6-41cf-a541-84efa50397c3?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=ciyXwq2TpzVGXEcQAqSdew%2F20190529%2Fru-msk%2Fs3%2Faws4_request&X-Amz-Date=20190529T123909Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&X-Amz-Signature=55d828f49ad5ecdaf07d40fb0c55003f5d6616ee8a2c631ca8e769636031a5b8"
		if id, erro := userSQL.ID.Value(); erro == nil {
			user.ID = int(id.(int64))
		}
		if name, _ := userSQL.Name.Value(); name != nil {
			user.Name = name.(string)
		}
		if photoURL, _ := userSQL.PhotoURL.Value(); photoURL != nil {
			user.FileKey = photoURL.(string)
		}

		fmt.Println("load message:", user.Name, user.PhotoURL)

		messages = append(messages, message)
	}
	if err != nil {
		fmt.Println("database/GetUsers wrong row catched:", err.Error())
		return
	}

	fmt.Println("database/getMessages +")

	return
}
