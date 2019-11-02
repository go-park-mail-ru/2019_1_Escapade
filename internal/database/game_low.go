package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"database/sql"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) createGame(tx *sql.Tx, game *models.Game) (id int32, err error) {
	sqlInsert := `
	INSERT INTO Game(roomID, name, players, timeToPrepare,
		timeToPlay, date, noAnonymous, deathmatch) VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
		`
	row := tx.QueryRow(sqlInsert, game.Settings.ID, game.Settings.Name,
		game.Settings.Players, game.Settings.TimeToPrepare,
		game.Settings.TimeToPlay, game.Date, game.Settings.NoAnonymous,
		game.Settings.Deathmatch)

	err = row.Scan(&id)

	return
}

func (db *DataBase) updateGame(tx *sql.Tx, game models.Game) error {
	sqlStatement := `
	UPDATE Game 
		SET status = $1, chatID = $2, recruitment = $3, playing = $4
		WHERE id = $5
		RETURNING id
	`

	row := tx.QueryRow(sqlStatement, game.Status, game.ChatID,
		game.RecruitmentTime, game.PlayingTime, game.Settings.ID)

	var id int32
	err := row.Scan(&id)
	return err
}

func (db *DataBase) createGamers(tx *sql.Tx, GameID int32, gamers []models.Gamer) (err error) {
	sqlInsert := `
	INSERT INTO Gamer(player_id, game_id, score, time,
		 left_click, right_click, explosion, won) VALUES
    ($1, $2, $3, $4::interval, $5, $6, $7, $8);
		`

	for _, gamer := range gamers {
		_, err = tx.Exec(sqlInsert, gamer.ID, GameID, gamer.Score, gamer.Time,
			gamer.LeftClick, gamer.RightClick, gamer.Explosion, gamer.Won)
		if err != nil {
			break
		}
	}
	return
}

func (db *DataBase) createField(tx *sql.Tx, gameID int32, field models.Field) (id int32, err error) {
	sqlInsert := `
	INSERT INTO Field(game_id, width, height, cells_left, difficult,
		mines) VALUES
		($1, $2, $3, $4, $5, $6)
		RETURNING id
		`
	row := tx.QueryRow(sqlInsert, gameID, field.Width,
		field.Height, field.CellsLeft, field.Difficult,
		field.Mines)

	err = row.Scan(&id)
	return
}

func (db *DataBase) createActions(tx *sql.Tx, GameID int32, actions []models.Action) (err error) {
	sqlInsert := `
	INSERT INTO Action(game_id, player_id, action, date) VALUES
    ($1, $2, $3, $4);
		`

	for _, action := range actions {
		_, err = tx.Exec(sqlInsert, GameID,
			action.PlayerID, action.ActionID, action.Date)
		if err != nil {
			break
		}
	}
	return
}

func (db *DataBase) createCells(tx *sql.Tx, FieldID int32, cells []models.Cell) (err error) {
	sqlInsert := `
	INSERT INTO Cell(field_id, player_id, x, y, value, date) VALUES
    ($1, $2, $3, $4, $5, $6);
		`

	for _, cell := range cells {
		_, err = tx.Exec(sqlInsert, FieldID, cell.PlayerID,
			cell.X, cell.Y, cell.Value, cell.Date)
		if err != nil {
			break
		}
	}
	return
}

// getGamesURL get user games URL
func (db *DataBase) getGamesURL(tx *sql.Tx, playerID int32) (URLs []string, err error) {
	getURLs := `
	SELECT roomID
				 FROM Game
				 join Gamer
				 on Game.id = Gamer.game_id 
				 where player_id = $1
	`

	URLs = make([]string, 0)
	rows, erro := tx.Query(getURLs, playerID)

	if erro != nil {
		err = erro
		return
	}
	defer rows.Close()

	for rows.Next() {
		var url string
		if err = rows.Scan(&url); err != nil {

			break
		}
		URLs = append(URLs, url)
	}
	if err != nil {
		return
	}
	return
}

// GetGameInformation get all information about game:
// game, gamers, field, history of cells and actions
func (db *DataBase) GetGameInformation(tx *sql.Tx, roomID string) (gameInformation models.GameInformation, err error) {

	getGame := `
	SELECT id, roomID, name, players, status, timeToPrepare,
	 timeToPlay, date 
				 FROM Game
				 where roomID like $1
	`

	row := tx.QueryRow(getGame, roomID)

	game := models.Game{}
	var gameID int32
	err = row.Scan(&gameID, &game.Settings.ID, &game.Settings.Name,
		&game.Settings.Players, &game.Status, &game.Settings.TimeToPrepare,
		&game.Settings.TimeToPlay, &game.Date)
	if err != nil {
		return
	}

	getGamers := `
	SELECT GR.player_id, GR.score, EXTRACT(seconds FROM GR.time), GR.left_click,
				GR.right_click, GR.explosion, GR.won
			FROM Gamer as GR 
			where GR.game_id = $1
	`

	gamers := make([]models.Gamer, 0)
	rows, erro := tx.Query(getGamers, gameID)

	if erro != nil {
		err = erro
		return
	}
	defer rows.Close()

	for rows.Next() {
		gamer := models.Gamer{}
		if err = rows.Scan(&gamer.ID, &gamer.Score, &gamer.Time, &gamer.LeftClick,
			&gamer.RightClick, &gamer.Explosion, &gamer.Won); err != nil {

			break
		}
		gamers = append(gamers, gamer)
	}
	if err != nil {
		return
	}

	getField := `
	SELECT id, width, height, cells_left, difficult, mines
		from Field where game_id = $1
	`

	row = tx.QueryRow(getField, gameID)

	field := models.Field{}
	var fieldID int
	err = row.Scan(&fieldID, &field.Width, &field.Height,
		&field.CellsLeft, &field.Difficult, &field.Mines)
	if err != nil {
		return
	}

	getActions := `
	SELECT player_id, action, date
		from Action where game_id = $1
	`

	actions := make([]models.Action, 0)
	rows1, err1 := tx.Query(getActions, gameID)

	if err1 != nil {
		err = err1
		return
	}
	defer rows1.Close()

	for rows1.Next() {
		action := models.Action{}
		if err = rows1.Scan(&action.PlayerID, &action.ActionID,
			&action.Date); err != nil {

			break
		}
		actions = append(actions, action)
	}
	if err != nil {
		return
	}

	getCells := `
	SELECT player_id, x, y, value, date
		from Cell where field_id = $1
	`

	cells := make([]models.Cell, 0)
	rows2, err2 := tx.Query(getCells, fieldID)

	if err2 != nil {
		err = err2
		return
	}
	defer rows2.Close()

	for rows2.Next() {
		cell := models.Cell{}
		if err = rows2.Scan(&cell.PlayerID, &cell.X,
			&cell.Y, &cell.Value, &cell.Date); err != nil {

			break
		}
		cells = append(cells, cell)
	}
	if err != nil {
		return
	}
	/*

		var (
			chat = &pChat.Chat{
				Type:   pChat.ChatType_ROOM,
				TypeId: gameID,
			}
			chatID    *pChat.ChatID
			pMessages *pChat.Messages
		)

		chatID, err = clients.ALL.Chat.GetChat(context.Background(), chat)

		if err != nil {
			utils.Debug(true, "cant access to chat service", err.Error())
		}
		pMessages, err = clients.ALL.Chat.GetChatMessages(context.Background(), chatID)

		var messages []*models.Message
		messages = MessagesFromProto(pMessages.Messages...)
		//db.getMessages(tx, true, game.RoomID)
		if err != nil {
			utils.Debug(true, "cant get messages!", err.Error())
		}
	*/

	return models.GameInformation{
		Game:    game,
		Gamers:  gamers,
		Field:   field,
		Actions: actions,
		Cells:   cells,
		//Messages: messages,
	}, err

}
