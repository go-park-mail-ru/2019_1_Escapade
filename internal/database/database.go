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
	Db *sql.DB
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
	sqlStatement := `UPDATE Player SET photo = $1 WHERE name = $2;`

	_, err = db.Db.Exec(sqlStatement, filename, username)

	if err != nil {
		fmt.Println("database/session/PostImage - fail:" + err.Error())
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

func (db *DataBase) GetUsers(name string, how int) (games []models.Game, err error) {

	sqlStatement := `
	SELECT * 
	FROM Player as P1
	JOIN (
		SELECT email, best_score, best_time  
		FROM Player 
		ORDER BY id LIMIT 100000 OFFSET 2)
		as P2 ON b.id = test_table.id
	SELECT SELECT email, best_score, best_time 
	FROM Player 
	ORDER BY (best_score)
`
	games = make([]models.Game, 0, 0)
	rows, erro := db.Db.Query(sqlStatement, name)

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
func (db *DataBase) GetGames(name string) (games []models.Game, err error) {

	sqlStatement := `
	SELECT FieldWidth, FieldHeight, MinsTotal, MinsFound,
				 Finished, Exploded
	FROM Game as G join Player as P on G.player_id=P.id
	WHERE P.name like $1 
`
	games = make([]models.Game, 0, 0)
	rows, erro := db.Db.Query(sqlStatement, name)

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
