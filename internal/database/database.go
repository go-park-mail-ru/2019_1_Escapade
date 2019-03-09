package database

import (
	"database/sql"
	"escapade/internal/models"

	"fmt"

	//
	_ "github.com/lib/pq"
)

// DataBase consists of *sql.DB
// Support methods Login, Register
type DataBase struct {
	Db        *sql.DB
	PageGames int
	PageUsers int
}

// Login check sql-injections and is password right
// Then add cookie to database and returns session_id
func (db *DataBase) Login(user *models.UserPrivateInfo) (str string, err error) {

	if err = ValidatePrivateUI(user); err != nil {
		fmt.Println("database/login - fail validation")
		return
	}

	if user.Name == "" {
		fmt.Println("+")
		if user.Name, err = GetNameByEmail(user.Email, db.Db); err != nil {
			fmt.Println("database/login - fail get name by email")
			return
		}
	}
	fmt.Println("User", user.Name, user.Email)
	if err = confirmRightPass(user, db.Db); err != nil {
		fmt.Println("database/login - fail confirmition")
		return
	}

	if err = confirmRightPass(user, db.Db); err != nil {
		fmt.Println("database/login - fail confirmition")
		return
	}

	if str, err = db.createSession(user); err != nil {
		fmt.Println("database/login - fail creating Session")
		return
	}

	fmt.Println("database/login +")

	return
}

// Register check sql-injections and are email and name unique
// Then add cookie to database and returns session_id
func (db *DataBase) Register(user *models.UserPrivateInfo) (str string, err error) {

	if err = ValidatePrivateUI(user); err != nil {
		fmt.Println("database/register - fail validation")
		return
	}

	if err = confirmUnique(user, db.Db); err != nil {
		fmt.Println("database/register - fail uniqie")
		return
	}

	if err = db.createUser(user); err != nil {
		fmt.Println("database/register - fail creating User")
		return
	}

	if str, err = db.createSession(user); err != nil {
		fmt.Println("database/register - fail creating Session")
		return
	}

	fmt.Println("database/register +")

	return
}

// Logout delete session_id row  from session table
func (db *DataBase) Logout(sessionCode string) (err error) {
	err = db.deleteSession(sessionCode)
	return
}

func (db *DataBase) PostImage(filename string, username string) (err error) {
	sqlStatement := `UPDATE Player SET photo_title = $1 WHERE name = $2;`

	_, err = db.Db.Exec(sqlStatement, filename, username)

	if err != nil {
		fmt.Println("database/session/PostImage - fail:" + err.Error())
		return
	}
	return
}

func (db *DataBase) GetImage(username string) (filename string, err error) {
	sqlStatement := `
	SELECT photo_title
		FROM Player as P 
			join Session as S 
			on S.player_id=P.id
		WHERE P.name = $1 
`
	row := db.Db.QueryRow(sqlStatement, username)

	if err = row.Scan(&filename); err != nil {
		fmt.Println("database/GetImage failed")
		return
	}
	return
}

func (db *DataBase) GetNameBySessionID(sessionID string) (name string, err error) {
	sqlStatement := `
	SELECT name
	FROM Player as P join Session as S on S.player_id=P.id
	WHERE session_code like $1 
`
	row := db.Db.QueryRow(sqlStatement, sessionID)

	err = row.Scan(&name)
	if err != nil {
		fmt.Println("Sess error: ", err.Error())
		fmt.Println("database/GetNameBySessionID failed")
		return
	}

	return
}

func (db *DataBase) GetUsersPageAmount() (amount int, err error) {
	sqlStatement := `SELECT count(1) FROM Player`
	row := db.Db.QueryRow(sqlStatement)
	if err = row.Scan(&amount); err != nil {
		fmt.Println("GetUsersAmount failed")
		return
	}
	amount /= db.PageUsers
	amount++
	fmt.Println("Page amount", amount)
	return
}

// GetUsers returns information about users
// for leaderboard
func (db *DataBase) GetUsers(page int) (players []models.UserPublicInfo, err error) {

	sqlStatement := `
	SELECT P1.name, P1.email, P1.best_score, P1.best_time  
	FROM Player as P1
	JOIN (
		SELECT id, name, email, best_score, best_time  
		FROM Player
		OFFSET $1 Limit $2
		)
		as P2 ON P1.id = P2.id
`
	size := db.PageUsers
	players = make([]models.UserPublicInfo, 0, size)
	rows, erro := db.Db.Query(sqlStatement, size*(page-1), size)

	if erro != nil {
		err = erro

		fmt.Println("database/GetUsers cant access to database:", erro.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		player := models.UserPublicInfo{}
		if err = rows.Scan(&player.Name, &player.Email, &player.BestScore,
			&player.BestTime); err != nil {

			fmt.Println("database/GetUsers wrong row catched")

			break
		}

		players = append(players, player)
	}

	fmt.Println("database/GetUsers +")

	return
}

func (db *DataBase) GetProfile(name string) (user models.UserPublicInfo, err error) {

	sqlStatement := `
	SELECT email, best_score, best_time 
	FROM Player 
	WHERE name like $1 
`

	row := db.Db.QueryRow(sqlStatement, name)

	user.Name = name

	if err = row.Scan(&user.Email, &user.BestScore, &user.BestTime); err != nil {
		fmt.Println("database/GetProfile failed")
		return
	}

	fmt.Println("database/GetProfile +")

	return
}

// GetGames returns games, played by player with some name
func (db *DataBase) GetGames(name string, page int) (games []models.Game, err error) {

	size := db.PageGames
	sqlStatement := `
	SELECT 	a.FieldWidth, a.FieldHeight,
					a.MinsTotal, a.MinsFound,
					a.Finished, a.Exploded 
	 FROM Player as p 
		JOIN
			(
				SELECT player_id,
					FieldWidth, FieldHeight,
					MinsTotal, MinsFound,
					Finished, Exploded 
					FROM Game Order by id
			) as a
			ON p.id = a.player_id and p.name like $1
			OFFSET $2 Limit $3
	`

	games = make([]models.Game, 0, size)
	rows, erro := db.Db.Query(sqlStatement, name, size*(page-1), size) // //, name)

	if erro != nil {
		err = erro

		fmt.Println("database/GetGames cant access to database:", erro.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		game := models.Game{}
		if err = rows.Scan(&game.FieldWidth, &game.FieldHeight,
			&game.MinsTotal, &game.MinsFound, &game.Finished,
			&game.Exploded); err != nil {

			fmt.Println("database/GetGames wrong row catched")

			break
		}

		games = append(games, game)
	}

	fmt.Println("database/GetGames +")

	return
}

// DeleteAccount deletes account
func (db *DataBase) DeleteAccount(user *models.UserPrivateInfo, sessionCode string) (str string, err error) {

	if err = ValidatePrivateUI(user); err != nil {
		fmt.Println("database/DeleteAccount - fail validation")
		return
	}

	if err = confirmRightPass(user, db.Db); err != nil {
		fmt.Println("database/DeleteAccount - fail confirmition password")
		return
	}

	if err = confirmRightEmail(user, db.Db); err != nil {
		fmt.Println("database/DeleteAccount - fail confirmition email")
		return
	}

	if err = db.deleteSession(sessionCode); err != nil {
		fmt.Println("database/DeleteAccount - fail deleting Session")
		return
	}

	if err = db.deleteUser(user); err != nil {
		fmt.Println("database/DeleteAccount - fail deletting User")
		return
	}

	fmt.Println("database/DeleteAccount +")

	return
}
