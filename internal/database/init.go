package database

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"

	"database/sql"
	"fmt"
	"os"

	//
	_ "github.com/lib/pq"
)

// Init try to connect to DataBase.
// If success - return instance of DataBase
// if failed - return error
func Init(CDB config.DatabaseConfig) (db *DataBase, err error) {

	var database *sql.DB
	if database, err = sql.Open(CDB.DriverName, os.Getenv(CDB.URL)); err != nil {
		fmt.Println("database/Init cant open:" + err.Error())
		return
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

	return
}

// InitWithRebuild connect to db with drop/create tables
func InitWithRebuild(CDB config.DatabaseConfig) (db *DataBase, err error) {

	if db, err = Init(CDB); err != nil {
		return
	}

	var (
		tx *sql.Tx
	)

	if tx, err = db.Db.Begin(); err != nil {
		return
	}
	defer tx.Rollback()

	if err = db.dropTables(tx); err != nil {
		return
	}

	if err = db.createTables(tx); err != nil {
		return
	}

	err = tx.Commit()

	return
}

func (db *DataBase) dropTables(tx *sql.Tx) (err error) {
	sqlStatement := `
    DROP TABLE IF EXISTS Session cascade;
    DROP TABLE IF EXISTS Cell cascade;
    DROP TABLE IF EXISTS Action cascade;
    DROP TABLE IF EXISTS Field cascade;
    DROP TABLE IF EXISTS Game cascade;
    DROP TABLE IF EXISTS GameChat cascade;
    DROP TABLE IF EXISTS Player cascade;
    DROP TABLE IF EXISTS Gamer cascade;
    DROP TABLE IF EXISTS Record cascade;
    `
	_, err = tx.Exec(sqlStatement)
	return err
}

func (db *DataBase) createTables(tx *sql.Tx) (err error) {
	if err = db.createTablePlayer(tx); err != nil {
		return
	}

	if err = db.createTableSession(tx); err != nil {
		return
	}

	if err = db.createTableRecord(tx); err != nil {
		return
	}

	if err = db.createTableGame(tx); err != nil {
		return
	}

	if err = db.createTableGameChat(tx); err != nil {
		return
	}

	if err = db.createTableGamer(tx); err != nil {
		return
	}

	if err = db.createTableField(tx); err != nil {
		return
	}

	if err = db.createTableCell(tx); err != nil {
		return
	}

	if err = db.createTableAction(tx); err != nil {
		return
	}
	return
}

func (db *DataBase) createTablePlayer(tx *sql.Tx) (err error) {
	sqlStatement := `
	CREATE TABLE Player (
        id SERIAL PRIMARY KEY,
        name varchar(30) NOT NULL,
        password varchar(30) NOT NULL,
		photo_title varchar(50) default '1.png',
        firstSeen   TIMESTAMPTZ,
        lastSeen    TIMESTAMPTZ
    );

    CREATE UNIQUE INDEX idx_lower_unique 
        ON Player (lower(name));
    `
	_, err = tx.Exec(sqlStatement)
	return err
}

func (db *DataBase) createTableSession(tx *sql.Tx) (err error) {
	sqlStatement := `
	CREATE Table Session (
        id SERIAL PRIMARY KEY,
        player_id int NOT NULL,
        session_code varchar(30) NOT NULL
    );
    
    ALTER TABLE Session
    ADD CONSTRAINT session_player
       FOREIGN KEY (player_id)
       REFERENCES Player(id)
       ON DELETE CASCADE;
    `
	_, err = tx.Exec(sqlStatement)
	return err
}

func (db *DataBase) createTableRecord(tx *sql.Tx) (err error) {
	sqlStatement := `
	CREATE Table Record (
        id SERIAL PRIMARY KEY,
        player_id int NOT NULL,
        score int default 0,
        time interval default '24 hour'::interval,
        difficult int default 0,
        singleTotal int default 0 CHECK (SingleTotal > -1),
        onlineTotal int default 0 CHECK (OnlineTotal > -1),
        singleWin   int default 0 CHECK (SingleWin > -1),
        onlineWin   int default 0 CHECK (OnlineWin > -1)
    );
    
    ALTER TABLE Record
    ADD CONSTRAINT record_player
       FOREIGN KEY (player_id)
       REFERENCES Player(id)
       ON DELETE CASCADE;
    `
	_, err = tx.Exec(sqlStatement)
	return err
}

func (db *DataBase) createTableGame(tx *sql.Tx) (err error) {
	sqlStatement := `
	CREATE Table Game (
        id SERIAL PRIMARY KEY,
        roomID varchar(30),
        name varchar(30),
        players   int NOT NULL,
        status int NOT NULL,
        timeToPrepare int,
        timeToPlay int,
        date TIMESTAMPTZ not null
    );
    `
	_, err = tx.Exec(sqlStatement)
	return err
}

func (db *DataBase) createTableGameChat(tx *sql.Tx) (err error) {
	sqlStatement := `
	CREATE Table GameChat (
        id SERIAL PRIMARY KEY,
        in_room bool,
        name varchar(30),
        roomID varchar(20),
        player_id int NOT NULL,
        message varchar(8000),
        time   TIMESTAMPTZ,
        edited bool default false
    );
    `
	_, err = tx.Exec(sqlStatement)
	return err
}

func (db *DataBase) createTableField(tx *sql.Tx) (err error) {
	sqlStatement := `
	CREATE Table Field (
        id SERIAL PRIMARY KEY,
        game_id int NOT NULL,
        width   int NOT NULL,
        height   int NOT NULL,
        cellsLeft int NOT NULL,
        difficult int default 0,
        mines   int NOT NULL
    );

    ALTER TABLE Field
        ADD CONSTRAINT field_game
           FOREIGN KEY (game_id)
           REFERENCES Game(id)
           ON DELETE CASCADE;
    
    `
	_, err = tx.Exec(sqlStatement)
	return err
}

func (db *DataBase) createTableAction(tx *sql.Tx) (err error) {
	sqlStatement := `
	CREATE Table Action (
        id SERIAL PRIMARY KEY,
        game_id int NOT NULL,
        player_id   int NOT NULL,
        action int  NOT NULL,
        date TIMESTAMPTZ not null
    );

    ALTER TABLE Action
        ADD CONSTRAINT action_game
           FOREIGN KEY (game_id)
           REFERENCES Game(id)
           ON DELETE CASCADE;
/*
    ALTER TABLE Action
        ADD CONSTRAINT action_player
            FOREIGN KEY (player_id)
            REFERENCES Player(id)
            ON DELETE CASCADE;
            */
    
    `
	_, err = tx.Exec(sqlStatement)
	return err
}

/*
difficult int default 0,
        width   int NOT NULL,
        height   int NOT NULL,
        players   int NOT NULL,
        mines   int NOT NULL,
*/

func (db *DataBase) createTableGamer(tx *sql.Tx) (err error) {
	sqlStatement := `
	CREATE Table Gamer (
        id SERIAL PRIMARY KEY,
        player_id int NOT NULL,
        game_id int NOT NULL,
        score float default 0,
        time interval default '24 hour'::interval,
        left_click int default 0,
        right_click int default 0,
        explosion bool default false,
        won bool default false
    );
    /*
    ALTER TABLE Gamer
    ADD CONSTRAINT gamer_player
       FOREIGN KEY (player_id)
       REFERENCES Player(id)
       ON DELETE CASCADE;
    */
    
    ALTER TABLE Gamer
    ADD CONSTRAINT gamer_game
        FOREIGN KEY (game_id)
        REFERENCES Game(id)
        ON DELETE CASCADE;
    `
	_, err = tx.Exec(sqlStatement)
	return err
}

func (db *DataBase) createTableCell(tx *sql.Tx) (err error) {
	// player_id maybe 0. It is mean that it is set by room
	// thats why there is no constraint with player_id
	sqlStatement := `
	CREATE Table Cell (
        id SERIAL PRIMARY KEY,
        field_id int NOT NULL,
        player_id int,
        x   int NOT NULL,
        y   int NOT NULL,
        value   int NOT NULL,
        date TIMESTAMPTZ not null
    );
    
    ALTER TABLE Cell
    ADD CONSTRAINT cell_field
        FOREIGN KEY (field_id)
        REFERENCES Field(id)
        ON DELETE CASCADE;
    `
	_, err = tx.Exec(sqlStatement)
	return err
}
