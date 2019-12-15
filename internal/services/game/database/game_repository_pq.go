package database

import (
	idb "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/models"
)

// GameRepositoryPQ implements the interface GameRepositoryI using the sql postgres driver
type GameRepositoryPQ struct{}

// CreateGame create game
func (db *GameRepositoryPQ) CreateGame(tx idb.TransactionI, game *models.Game) (int32, error) {
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

	var id int32
	err := row.Scan(&id)
	return id, err
}

// UpdateGame update game
func (db *GameRepositoryPQ) UpdateGame(tx idb.TransactionI, game *models.Game) error {
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

// CreateGamers create gamers
func (db GameRepositoryPQ) CreateGamers(tx idb.TransactionI, GameID int32,
	gamers []models.Gamer) error {
	sqlInsert := `
	INSERT INTO Gamer(player_id, game_id, score, time,
		 left_click, right_click, explosion, won) VALUES
    ($1, $2, $3, $4::interval, $5, $6, $7, $8);
		`
	var err error

	for _, gamer := range gamers {
		_, err = tx.Exec(sqlInsert, gamer.ID, GameID, gamer.Score, gamer.Time,
			gamer.LeftClick, gamer.RightClick, gamer.Explosion, gamer.Won)
		if err != nil {
			break
		}
	}
	return err
}

// CreateField create field
func (db GameRepositoryPQ) CreateField(tx idb.TransactionI, gameID int32,
	field models.Field) (int32, error) {
	sqlInsert := `
	INSERT INTO Field(game_id, width, height, cells_left, difficult,
		mines) VALUES
		($1, $2, $3, $4, $5, $6)
		RETURNING id
		`
	row := tx.QueryRow(sqlInsert, gameID, field.Width,
		field.Height, field.CellsLeft, field.Difficult,
		field.Mines)

	var id int32
	err := row.Scan(&id)
	return id, err
}

// CreateActions create actions
func (db GameRepositoryPQ) CreateActions(tx idb.TransactionI, GameID int32, actions []models.Action) error {
	sqlInsert := `
	INSERT INTO Action(game_id, player_id, action, date) VALUES
    ($1, $2, $3, $4);
		`

	var err error
	for _, action := range actions {
		_, err = tx.Exec(sqlInsert, GameID,
			action.PlayerID, action.ActionID, action.Date)
		if err != nil {
			break
		}
	}
	return err
}

// CreateCells create cells
func (db GameRepositoryPQ) CreateCells(tx idb.TransactionI, FieldID int32, cells []models.Cell) error {
	sqlInsert := `
	INSERT INTO Cell(field_id, player_id, x, y, value, date) VALUES
    ($1, $2, $3, $4, $5, $6);
		`

	var err error
	for _, cell := range cells {
		_, err = tx.Exec(sqlInsert, FieldID, cell.PlayerID,
			cell.X, cell.Y, cell.Value, cell.Date)
		if err != nil {
			break
		}
	}
	return err
}

// FetchOneGame fetch one game
func (db *GameRepositoryPQ) FetchOneGame(tx idb.TransactionI, roomID string) (models.Game, error) {

	getGame := `
	SELECT id, roomID, name, players, status, timeToPrepare, timeToPlay, date 
				 FROM Game
				 where roomID like $1
	`

	row := tx.QueryRow(getGame, roomID)

	game := models.Game{}
	err := row.Scan(&game.ID, &game.Settings.ID, &game.Settings.Name,
		&game.Settings.Players, &game.Status, &game.Settings.TimeToPrepare,
		&game.Settings.TimeToPlay, &game.Date)
	return game, err
}

// FetchAllRoomsID get user games URL
func (db *GameRepositoryPQ) FetchAllRoomsID(tx idb.TransactionI, playerID int32) ([]string, error) {
	var (
		getURLs = `SELECT roomID
				 FROM Game
				 join Gamer
				 on Game.id = Gamer.game_id 
				 where player_id = $1
	`

		URLs      = make([]string, 0)
		rows, err = tx.Query(getURLs, playerID)
	)

	if err != nil {
		return URLs, err
	}
	defer rows.Close()

	for rows.Next() {
		var url string
		if err = rows.Scan(&url); err != nil {
			break
		}
		URLs = append(URLs, url)
	}
	return URLs, err
}

// FetchAllGamers fetch all gamers
func (db *GameRepositoryPQ) FetchAllGamers(tx idb.TransactionI, gameID int32) ([]models.Gamer, error) {
	getGamers := `
	SELECT GR.player_id, GR.score, EXTRACT(seconds FROM GR.time), GR.left_click,
				GR.right_click, GR.explosion, GR.won
			FROM Gamer as GR 
			where GR.game_id = $1`

	gamers := make([]models.Gamer, 0)
	rows, err := tx.Query(getGamers, gameID)
	if err != nil {
		return gamers, err
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
	return gamers, err
}

func (db *GameRepositoryPQ) FetchOneField(tx idb.TransactionI, gameID int32) (int, models.Field, error) {
	getField := `SELECT id, width, height, cells_left, difficult, mines
		from Field where game_id = $1`
	row := tx.QueryRow(getField, gameID)
	var (
		field   models.Field
		fieldID int
	)
	err := row.Scan(&fieldID, &field.Width, &field.Height,
		&field.CellsLeft, &field.Difficult, &field.Mines)
	return fieldID, field, err
}

func (db *GameRepositoryPQ) FetchAllActions(tx idb.TransactionI, gameID int32) ([]models.Action, error) {
	getActions := ` SELECT player_id, action, date 
	from Action where game_id = $1`

	actions := make([]models.Action, 0)
	rows, err := tx.Query(getActions, gameID)
	if err != nil {
		return actions, err
	}
	defer rows.Close()

	for rows.Next() {
		action := models.Action{}
		err = rows.Scan(&action.PlayerID, &action.ActionID, &action.Date)
		if err != nil {
			break
		}
		actions = append(actions, action)
	}
	return actions, err
}

func (db *GameRepositoryPQ) FetchAllCells(tx idb.TransactionI, fieldID int) ([]models.Cell, error) {
	getCells := `SELECT player_id, x, y, value, date
from Cell where field_id = $1`

	cells := make([]models.Cell, 0)
	rows, err := tx.Query(getCells, fieldID)

	if err != nil {
		return cells, err
	}
	defer rows.Close()

	for rows.Next() {
		cell := models.Cell{}
		if err = rows.Scan(&cell.PlayerID, &cell.X,
			&cell.Y, &cell.Value, &cell.Date); err != nil {

			break
		}
		cells = append(cells, cell)
	}
	return cells, err
}
