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

	// for local launch
	if os.Getenv(CDB.URL) == "" {
		//db://postgres:postgres@db:5432/postgres?sslmode=disable
		//os.Setenv(CDB.URL, "postgresql://rolepade:escapade@localhost:5432/escabase")
		os.Setenv(CDB.URL, "user=docker password=docker dbname=docker sslmode=disable")
	}
	//os.Setenv(CDB.URL, "postgresql://rolepade:escapade@127.0.0.1:5432/escabase")
	fmt.Println("url:" + string(os.Getenv(CDB.URL)))

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
	//if !db.areTablesCreated(CDB.Tables) {
	if err = db.CreateTables(); err != nil {
		return
	}
	//}

	return
}

func (db *DataBase) checkTable(tableName string) (err error) {
	sqlStatement := `
    SELECT count(1)
  FROM information_schema.tables tbl 
  where tbl.table_name like $1;`
	row := db.Db.QueryRow(sqlStatement, tableName)

	var result int
	if err = row.Scan(&result); err != nil {
		fmt.Println(tableName + " doesnt exists. Create it!" + err.Error())

		return
	}
	return
}

func (db *DataBase) areTablesCreated(tables []string) (created bool) {
	created = true
	for _, table := range tables {
		if err := db.checkTable(table); err != nil {
			created = false
			break
		}
	}
	return
}

func (db *DataBase) CreateTables() error {
	sqlStatement := `
	DROP TABLE IF EXISTS Session;
    DROP TABLE IF EXISTS Game;
    DROP TABLE IF EXISTS Player;
    DROP TABLE IF EXISTS Photo;

	CREATE TABLE Player (
    id SERIAL PRIMARY KEY,
    name varchar(30) NOT NULL,
    password varchar(30) NOT NULL,
    email varchar(30) NOT NULL,
    photo_title varchar(50) default 'default',
    FirstSeen   TIMESTAMPTZ,
	LastSeen    TIMESTAMPTZ,
    best_score  int default 0 CHECK (best_score > -1),
    best_time   int default 0 CHECK (best_time > -1),
    GamesTotal  int default 0 CHECK (GamesTotal > -1),
	SingleTotal int default 0 CHECK (SingleTotal > -1),
	OnlineTotal int default 0 CHECK (OnlineTotal > -1),
	SingleWin   int default 0 CHECK (SingleWin > -1),
	OnlineWin   int default 0 CHECK (OnlineWin > -1),
	MinsFound   int default 0 CHECK (MinsFound > -1)
    
);

CREATE Table Session (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL,
    session_code varchar(30) NOT NULL,
    expiration timestamp without time zone NOT NULL
);

ALTER TABLE Session
ADD CONSTRAINT session_player
   FOREIGN KEY (player_id)
   REFERENCES Player(id)
   ON DELETE CASCADE;

CREATE Table Game (
    id SERIAL PRIMARY KEY,
    player_id   int NOT NULL,
    FieldWidth  int CHECK (FieldWidth > -1),
    FieldHeight int CHECK (FieldHeight > -1),
    MinsTotal   int CHECK (MinsTotal > -1),
    MinsFound   int CHECK (MinsFound > -1),
    Finished bool NOT NULL,
    Exploded bool NOT NULL,
    Date timestamp without time zone NOT NULL,
    FOREIGN KEY (player_id) REFERENCES Player (id)
);

--GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO escapade;

INSERT INTO Player(name, password, email, best_score, best_time) VALUES
    ('tiger', 'Bananas', 'tinan@mail.ru', 1000, 10),
    ('panda', 'apple', 'today@mail.ru', 2323, 20),
    ('catmate', 'juice', 'allday@mail.ru', 10000, 5),
    ('hotdog', 'where', 'three@mail.ru', 88, 1000),
    ('coala', 'cheese', 'thesame@mail.ru', 12050, 10),
    ('prten', 'apple', 'knowit@mail.ru', 23, 20),
    ('kingdom', 'notnot123', 'king@mail.ru', 10, 5),
    ('lifeIsStrange', 'always', 'life@mail.ru', 111, 1000),
    ('coala1', 'cheese', 'thesame@mail.ru', 12050, 10),
    ('prten1', 'apple', 'knowit@mail.ru', 23, 20),
    ('april', 'aprilapril', 'april@mail.ru', 10, 5),
    ('useruser', 'password', 'mail@mail.ru', 111, 1000),
    ('tiger1', 'Bananas', 'tinan@mail.ru', 1000, 10),
    ('panda1', 'apple', 'today@mail.ru', 2323, 20),
    ('catmate1', 'juice', 'allday@mail.ru', 10000, 5),
    ('hotdog1', 'where', 'three@mail.ru', 88, 1000),
    ('coala1', 'cheese', 'thesame@mail.ru', 12050, 10),
    ('prten1', 'apple', 'knowit@mail.ru', 23, 20),
    ('kingdom1', 'notnot123', 'king@mail.ru', 10, 5),
    ('lifeIsStrange1', 'always', 'life@mail.ru', 111, 1000),
    ('coala11', 'cheese', 'thesame@mail.ru', 12050, 10),
    ('prten11', 'apple', 'knowit@mail.ru', 23, 20),
    ('april1', 'aprilapril', 'april@mail.ru', 10, 5),
    ('useruser1', 'password', 'mail@mail.ru', 111, 1000),
    ('test1', 'password', 'mail@mail.ru', 0, 0),
    ('test2', 'password', 'mail@mail.ru', 0, 0),
    ('test3', 'password', 'mail@mail.ru', 0, 0),
    ('test4', 'password', 'mail@mail.ru', 0, 0),
    ('test5', 'password', 'mail@mail.ru', 0, 0),
    ('test6', 'password', 'mail@mail.ru', 0, 0),
    ('test7', 'password', 'mail@mail.ru', 0, 0),
    ('test8', 'password', 'mail@mail.ru', 0, 0),
    ('test9', 'password', 'mail@mail.ru', 0, 0),
    ('test10', 'password', 'mail@mail.ru', 0, 0),
    ('test11', 'password', 'mail@mail.ru', 0, 0),
    ('test12', 'password', 'mail@mail.ru', 0, 0),
    ('test13', 'password', 'mail@mail.ru', 0, 0),
    ('test14', 'password', 'mail@mail.ru', 0, 0),
    ('test15', 'password', 'mail@mail.ru', 0, 0),
    ('test16', 'password', 'mail@mail.ru', 0, 0),
    ('test17', 'password', 'mail@mail.ru', 0, 0),
    ('test18', 'password', 'mail@mail.ru', 0, 0),
    ('test19', 'password', 'mail@mail.ru', 0, 0),
    ('test20', 'password', 'mail@mail.ru', 0, 0),
    ('test21', 'password', 'mail@mail.ru', 0, 0),
    ('test22', 'password', 'mail@mail.ru', 0, 0),
    ('test23', 'password', 'mail@mail.ru', 0, 0),
    ('test24', 'password', 'mail@mail.ru', 0, 0),
    ('test25', 'password', 'mail@mail.ru', 0, 0),
    ('test26', 'password', 'mail@mail.ru', 0, 0),
    ('test27', 'password', 'mail@mail.ru', 0, 0),
    ('test28', 'password', 'mail@mail.ru', 0, 0),
    ('test29', 'password', 'mail@mail.ru', 0, 0),
    ('test30', 'password', 'mail@mail.ru', 0, 0),
    ('test40', 'password', 'mail@mail.ru', 0, 0),
    ('test41', 'password', 'mail@mail.ru', 0, 0),
    ('test42', 'password', 'mail@mail.ru', 0, 0),
    ('test43', 'password', 'mail@mail.ru', 0, 0),
    ('test44', 'password', 'mail@mail.ru', 0, 0),
    ('test45', 'password', 'mail@mail.ru', 0, 0),
    ('test46', 'password', 'mail@mail.ru', 0, 0),
    ('test47', 'password', 'mail@mail.ru', 0, 0),
    ('test48', 'password', 'mail@mail.ru', 0, 0),
    ('test49', 'password', 'mail@mail.ru', 0, 0),
    ('test50', 'password', 'mail@mail.ru', 0, 0),
    ('test51', 'password', 'mail@mail.ru', 0, 0),
    ('test52', 'password', 'mail@mail.ru', 0, 0),
    ('test53', 'password', 'mail@mail.ru', 0, 0),
    ('test54', 'password', 'mail@mail.ru', 0, 0),
    ('test55', 'password', 'mail@mail.ru', 0, 0),
    ('test56', 'password', 'mail@mail.ru', 0, 0),
    ('test57', 'password', 'mail@mail.ru', 0, 0),
    ('test58', 'password', 'mail@mail.ru', 0, 0),
    ('test59', 'password', 'mail@mail.ru', 0, 0),
    ('test60', 'password', 'mail@mail.ru', 0, 0);
    /*
    id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    name varchar(30) NOT NULL,
    password varchar(30) NOT NULL,
    email varchar(30) NOT NULL,
    photo_id int,
    best_score int,
    FOREIGN KEY (photo_id) REFERENCES Photo (id)
    */

INSERT INTO Game(player_id, FieldWidth, FieldHeight,
MinsTotal, MinsFound, Finished, Exploded, Date) VALUES
    (1, 50, 50, 100, 20, true, true, date '2001-09-28'),
    (1, 50, 50, 80, 30, false, false, date '2018-09-27'),
    (1, 50, 50, 70, 70, true, false, date '2018-09-26'),
    (1, 50, 50, 60, 30, true, true, date '2018-09-23'),
    (1, 50, 50, 50, 50, true, false, date '2018-09-24'),
    (1, 50, 50, 40, 30, true, false, date '2018-09-25'),
    (2, 25, 25, 80, 30, false, false, date '2018-08-27'),
    (2, 25, 25, 70, 70, true, false, date '2018-08-26'),
    (2, 25, 25, 60, 30, true, true, date '2018-08-23'),
    (3, 10, 10, 10, 10, true, false, date '2018-10-26'),
    (3, 10, 10, 20, 19, true, true, date '2018-10-23'),
    (3, 10, 10, 30, 30, true, false, date '2018-10-24'),
    (3, 10, 10, 40, 5, true, false, date '2018-10-25');

    /*
CREATE Table Game (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL,
    FieldWidth int CHECK (FieldWidth > -1),
    FieldHeight int CHECK (FieldHeight > -1),
    MinsTotal int CHECK (MinsTotal > -1),
    MinsFound int CHECK (MinsFound > -1),
    Finished bool NOT NULL,
    Exploded bool NOT NULL,
    Date timestamp without time zone NOT NULL,
    FOREIGN KEY (player_id) REFERENCES Player (id)
);
    */
	`
	_, err := db.Db.Exec(sqlStatement)

	if err != nil {
		fmt.Println("database/init - fail:" + err.Error())
	}
	return err
}
