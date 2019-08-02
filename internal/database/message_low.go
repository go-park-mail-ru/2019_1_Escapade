package database

/*
import (
	"fmt"
	"time"

	pChat "github.com/go-park-mail-ru/2019_1_Escapade/chat_service/proto"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"

	"strconv"
)*/

/*
import (
	"database/sql"

	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"
)

// GetUsers returns information about users
// for leaderboard
func (db *DataBase) getMessages(tx *sql.Tx, inRoom bool, gameID string) (messages []*models.Message, err error) {

	var (
		rows  *sql.Rows
		place = "getMessages -"
	)
	sqlStatement := `
	select GC.id, GC.player_id, GC.name, P.name, P.photo_title, GC.message, GC.time, GC.edited
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
		utils.Debug(true, place, "cant access to database:", err.Error())
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

		if err = rows.Scan(&message.ID, &userSQL.ID, &user.Name, &userSQL.Name,
			&userSQL.PhotoURL, &message.Text, &message.Time, &message.Edited); err != nil {

			break
		}
		user.FileKey = "anonymous.jpg"
		if id, erro := userSQL.ID.Value(); erro == nil {
			user.ID = int(id.(int64))
		}
		if name, _ := userSQL.Name.Value(); name != nil {
			user.Name = name.(string)
		}
		if photoURL, _ := userSQL.PhotoURL.Value(); photoURL != nil {
			user.FileKey = photoURL.(string)
		}

		messages = append(messages, message)
	}
	if err != nil {
		utils.Debug(true, place, "wrong row catched:", err.Error())
		return
	}

	return
}
*/
