package database

import (
	"database/sql"
	"escapade/internal/models"
	"os"

	"fmt"

	//
	_ "github.com/lib/pq"
)

// DataBase consists of *sql.DB
// Support methods Login, Register
type DataBase struct {
	Db *sql.DB
}

// Init try to connect to DataBase.
// If success - return instance of DataBase
// if failed - return error
func Init() (db *DataBase, err error) {
	//connStr := "user=unemuzhregdywt password=5d9ae1059a39b0a8838b5653854adc7fb266deb7da1dc35de729a4836ba27b65 dbname=dd1f3dqgsuq1k5 sslmode=disable"

	// connStr := "user=rolepade password=escapade dbname=escabase sslmode=disable"

	// var database *sql.DB
	// database, err = sql.Open("postgres", connStr)
	var database *sql.DB
	database, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return
	}
	db = &DataBase{Db: database}
	db.Db.SetMaxOpenConns(20)

	err = db.Db.Ping()
	if err != nil {
		return
	}

	return
}

// Login check sql-injections and is password right
// Then add cookie to database and returns session_id
func (db *DataBase) Login(user *models.UserPrivateInfo) (str string, err error) {

	if err = ValidatePrivateUI(user); err != nil {
		fmt.Println("database/login - fail validation")
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
