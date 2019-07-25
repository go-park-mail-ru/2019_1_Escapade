package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"database/sql"
	"os"

	//
	_ "github.com/lib/pq"
)

// Init try to connect to DataBase.
// If success - return instance of DataBase
// if failed - return error
func Init(CDB config.DatabaseConfig) (db *DataBase, err error) {

	var (
		database *sql.DB
		place    = "database Init() -"
	)
	if database, err = sql.Open(CDB.DriverName, os.Getenv(CDB.URL)); err != nil {
		utils.Debug(true, place, "cant open: -", err.Error())
		return
	}

	db = &DataBase{
		Db:        database,
		PageGames: CDB.PageGames,
		PageUsers: CDB.PageUsers,
	}
	db.Db.SetMaxOpenConns(CDB.MaxOpenConns)

	if err = db.Db.Ping(); err != nil {
		utils.Debug(true, place, "cant access: -", err.Error())
		return
	}
	utils.Debug(false, place, "success!")

	return
}
