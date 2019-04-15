package database

import (
	"crypto/rand"
	"database/sql"
	"escapade/internal/config"
	"escapade/internal/models"
	"fmt"
	ran "math/rand"
	"os"
	"time"

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
		os.Setenv(CDB.URL, "dbname=escabase user=rolepade password=escapade sslmode=disable")
		//"user=docker password=docker dbname=docker sslmode=disable")
	}

	os.Setenv("AWS_ACCESS_KEY_ID", "ciyXwq2TpzVGXEcQAqSdew")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "NzvtJoAid7GeUU2msVBzJXZGoA7rkjnQvnnEYZzujTx")

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
	// добавить в json проверку последнего апдейта
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
	DROP TABLE IF EXISTS Session cascade;
    DROP TABLE IF EXISTS Game cascade;
    DROP TABLE IF EXISTS Player cascade;
    DROP TABLE IF EXISTS Gamer cascade;
    DROP TABLE IF EXISTS Cell cascade;
    DROP TABLE IF EXISTS Record cascade;

	CREATE TABLE Player (
        id SERIAL PRIMARY KEY,
        name varchar(30) NOT NULL,
        password varchar(30) NOT NULL,
        email varchar(30) NOT NULL,
		photo_title varchar(50) default '1.png',
        firstSeen   TIMESTAMPTZ,
        lastSeen    TIMESTAMPTZ
    );

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

--GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO escapade;

INSERT INTO Player(name, password, email) VALUES
    ('tiger', 'Bananas', 'tinan@mail.ru'),
    ('panda', 'apple', 'today@mail.ru');

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

CREATE Table Game (
    id SERIAL PRIMARY KEY,
    difficult int default 0,
    width   int NOT NULL,
    height   int NOT NULL,
    players   int NOT NULL,
    mines   int NOT NULL,
    date TIMESTAMPTZ not null
);

CREATE Table Gamer (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL,
    game_id int NOT NULL,
    score int default 0,
    time interval default '24 hour'::interval,
    left_click int default 0,
    right_click int default 0,
    explosion bool default false,
    won bool default false
);

ALTER TABLE Gamer
ADD CONSTRAINT gamer_player
   FOREIGN KEY (player_id)
   REFERENCES Player(id)
   ON DELETE CASCADE;

ALTER TABLE Gamer
ADD CONSTRAINT gamer_game
    FOREIGN KEY (game_id)
    REFERENCES Game(id)
    ON DELETE CASCADE;

CREATE Table Cell (
    id SERIAL PRIMARY KEY,
    game_id int NOT NULL,
    gamer_id int NOT NULL,
    x   int NOT NULL,
    y   int NOT NULL,
    value   int NOT NULL
);

ALTER TABLE Cell
ADD CONSTRAINT cell_game
    FOREIGN KEY (game_id)
    REFERENCES Game(id)
    ON DELETE CASCADE;

ALTER TABLE Cell
ADD CONSTRAINT cell_gamer
    FOREIGN KEY (gamer_id)
    REFERENCES Gamer(id)
    ON DELETE CASCADE;

	`
	_, err := db.Db.Exec(sqlStatement)

	db.insert(110)

	if err != nil {
		fmt.Println("database/init - fail:" + err.Error())
	}
	return err
}

func (db *DataBase) insert(limit int) {

	n := 16
	for i := 0; i < limit; i++ {
		ran.Seed(time.Now().UnixNano())
		user := &models.UserPrivateInfo{
			Name:     RandString(n),
			Email:    RandString(n),
			Password: RandString(n)}
		_, id, _ := db.Register(user)
		for j := 0; j < 4; j++ {
			record := &models.Record{
				Score:       ran.Intn(1000000),
				Time:        float64(ran.Intn(10000)),
				Difficult:   j,
				SingleTotal: ran.Intn(2),
				OnlineTotal: ran.Intn(2),
				SingleWin:   ran.Intn(2),
				OnlineWin:   ran.Intn(2)}
			//fmt.Println("record:",record.Score, record.Time)
			db.UpdateRecords(id, record)
		}

	}
}

// RandString create random string with n length
func RandString(n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
