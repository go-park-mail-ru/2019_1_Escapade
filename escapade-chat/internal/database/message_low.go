package database

import (
	"database/sql"
	"escapade/internal/models"
	"fmt"
)

// createPlayer create player
func (db *DataBase) createMessage(tx *sql.Tx, message *models.Message) (id int, err error) {
	sqlInsert := `
	INSERT INTO UserChat(userID, name, photoUrl, message, time) VALUES
    ($1, $2, $3, $4, $5);
		`
	_, err = db.Db.Exec(sqlInsert, message.User.ID, message.User.Name,
		message.User.PhotoURL, message.Message, message.Time)

	if err != nil {
		fmt.Println("createMessage err:", err.Error())
		return
	}
	fmt.Println("createMessage success")

	return
}

// GetUsers returns information about users
// for leaderboard
func (db *DataBase) getMessages(tx *sql.Tx) (messages []*models.Message, err error) {

	sqlStatement := `select userID, name, photoUrl, message, time from UserChat;`

	messages = make([]*models.Message, 0)
	rows, erro := tx.Query(sqlStatement)

	if erro != nil {
		err = erro

		fmt.Println("database/getMessages cant access to database:", erro.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		user := &models.UserPublicInfo{}
		message := &models.Message{
			User: user,
		}
		if err = rows.Scan(&message.User.ID, &message.User.Name,
			&message.User.PhotoURL, &message.Message, &message.Time); err != nil {

			fmt.Println("database/GetUsers wrong row catched")

			break
		}
		fmt.Println("message:", message)

		messages = append(messages, message)
	}

	fmt.Println("database/getMessages +")

	return
}

// GetUsers returns information about users
// for leaderboard
func (db *DataBase) getUser(tx *sql.Tx, userID int, difficult int) (player *models.UserPublicInfo, err error) {

	sqlStatement := `
	SELECT P.id, P.photo_title, P.name, P.email,
				 R.score, R.time, R.Difficult
	FROM Player as P
	join Record as R 
	on R.player_id = P.id
	where R.player_id = $1 and
		R.difficult = $2
	`

	player = &models.UserPublicInfo{}
	row := tx.QueryRow(sqlStatement, userID, difficult)
	err = row.Scan(&player.ID, &player.FileKey, &player.Name,
		&player.Email, &player.BestScore, &player.BestTime, &player.Difficult)
	return
}

func (db *DataBase) deletePlayer(tx *sql.Tx, user *models.UserPrivateInfo) error {
	sqlStatement := `
	DELETE FROM Player where name=$1 and password=$2 and email=$3
		`
	fmt.Println("+++++")

	fmt.Println("user.Name, user.Password, user.Email", user.Name, user.Password, user.Email)

	_, err := tx.Exec(sqlStatement, user.Name,
		user.Password, user.Email)

	return err
}
