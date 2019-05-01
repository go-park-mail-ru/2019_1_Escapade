package database

import (
	"database/sql"
	"escapade/internal/config"
	"fmt"
	"os"

	//
	_ "github.com/lib/pq"
)

// Init try to connect to DataBase.
// If success - return instance of DataBase
// if failed - return error
func Init(CDB config.DatabaseConfig) (db *DataBase, err error) {
	//"postgres://docker:docker@tcp(db:5432)/docker"
	// for local launch
	if os.Getenv(CDB.URL) == "" {
		os.Setenv(CDB.URL, "dbname=docker host=localhost port=5432 user=docker password=docker sslmode=disable")
	}

	fmt.Println("url:" + string(os.Getenv(CDB.URL)))

	var database *sql.DB
	if database, err = sql.Open(CDB.DriverName, os.Getenv(CDB.URL)); err != nil {
		fmt.Println("database/Init cant open:" + err.Error())
	}

	db = &DataBase{
		Db:        database,
		PageGames: CDB.PageGames,
		PageUsers: CDB.PageUsers,
	}
	db.Db.SetMaxOpenConns(CDB.MaxOpenConns)

	if err = db.Db.Ping(); err != nil {
		fmt.Println("database/Init cant access:" + err.Error())
		return
	}
	fmt.Println("database/Init open")

	// добавить в json проверку последнего апдейта
	if err = db.CreateTables(); err != nil {
		return
	}

	return
}

// CreateTables drop old tables and create new
func (db *DataBase) CreateTables() error {
	sqlStatement := `
	DROP TABLE IF EXISTS UserChat cascade;

	CREATE TABLE UserChat (
		id SERIAL PRIMARY KEY,
				userID int NOT NULL,
				name varchar(30),
				photoUrl varchar(200),
				message varchar(8000),
        time   TIMESTAMPTZ
    );

	`
	_, err := db.Db.Exec(sqlStatement)

	if err != nil {
		fmt.Println("database/init - fail:" + err.Error())
	}
	return err
}
