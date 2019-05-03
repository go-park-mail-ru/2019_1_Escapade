package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"

	"database/sql"

	//
	_ "github.com/lib/pq"
)

func (db *DataBase) createGame(tx *sql.Tx, game models.Game) (id int, err error) {
	sqlInsert := `
	INSERT INTO Game(roomID, name, players, status, timeToPrepare,
		timeToPlay, date) VALUES
		($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
		`
	row := tx.QueryRow(sqlInsert, game.RoomID, game.Name,
		game.Players, game.Status, game.TimeToPrepare,
		game.TimeToPlay, game.Date)

	err = row.Scan(&id)

	return
}

func (db *DataBase) createGamers(tx *sql.Tx, GameID int, gamers []models.Gamer) (err error) {
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

func (db *DataBase) createField(tx *sql.Tx, gameID int, field models.Field) (id int, err error) {
	sqlInsert := `
	INSERT INTO Field(game_id, width, height, cellsLeft, difficult,
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

func (db *DataBase) createActions(tx *sql.Tx, GameID int, actions []models.Action) (err error) {
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

/*
type GameInformation struct {
	Game    Game     `json:"game"`
	Field   Field    `json:"field"`
	Actions []Action `json:"actions"`
	Cells   []Cell   `json:"cells"`
	Gamers  []Gamer  `json:"gamer"`
}
*/
func (db *DataBase) createCells(tx *sql.Tx, FieldID int, cells []models.Cell) (err error) {
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
func (db *DataBase) getGamesURL(tx *sql.Tx, playerID int) (URLs []string, err error) {
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
	var gameID int
	err = row.Scan(&gameID, &game.RoomID, &game.Name,
		&game.Players, &game.Status, &game.TimeToPrepare,
		&game.TimeToPlay, &game.Date)
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
	SELECT id, width, height, cellsLeft, difficult, mines
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

	return models.GameInformation{
		Game:    game,
		Gamers:  gamers,
		Field:   field,
		Actions: actions,
		Cells:   cells,
	}, err

}

// GetFullGamesInformation returns games, played by player with some name
/*
func (db *DataBase) GetFullGamesInformation(tx *sql.Tx, UserID int,
	page int) (games []*models.GameInformation, err error) {

	size := db.PageGames
	sqlStatement := `
	SELECT 	GE.width, GE.height, GE.difficult, GE.players,
	 GE.mines, GE.date, GR.score,
		GR.time, GR.left_click, GR.right_click, GR.explosion,
		GR.won
	 FROM Gamer as GR
	 join Game as GE
		ON GR.id = $1 and GR.game_id = GE.id
		OFFSET $2 Limit $3
	`

	games = make([]*models.GameInformation, 0, size)
	rows, erro := tx.Query(sqlStatement, UserID, size*(page-1), size) // //, name)

	if erro != nil {
		err = erro

		return
	}
	defer rows.Close()

	for rows.Next() {
		game := &models.GameInformation{}
		game.Game = &models.Game{}
		game.Gamers = make([]*models.Gamer, 1)
		/*
			ge.width, ge.height, ge.difficult,
							ge.players, ge.mines, ge.date, ge.online,
							gr.score, gr.time, gr.mines_open,
							gr.left_click, gr.right_click,
							gr.explosion, gr.won
			 FROM Player as p */
/*
		if err = rows.Scan(&game.Game.Width,
			&game.Game.Height,
			&game.Game.Difficult, &game.Game.Players,
			&game.Game.Mines, &game.Game.Date,
			&game.Gamers[0].Score, &game.Gamers[0].Time,
			&game.Gamers[0].LeftClick, &game.Gamers[0].RightClick,
			&game.Gamers[0].Explosion, &game.Gamers[0].Won); err != nil {

			break
		}

		games = append(games, game)
	}

	return
}*/
