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
		best_score int default 0,
		best_time int default 0,
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
/*
INSERT INTO Player(name, password, email, best_score, best_time) VALUES
    ('tiger', 'Bananas', 'tinan@mail.ru', 1000, '1 second'::interval),
    ('panda', 'apple', 'today@mail.ru', 2323, '2 second'::interval),
    ('catmate', 'juice', 'allday@mail.ru', 10000, '3 second'::interval),
    ('hotdog', 'where', 'three@mail.ru', 88, '4 second'::interval),
    ('coala', 'cheese', 'thesame@mail.ru', 12050, '5 second'::interval),
    ('prten', 'apple', 'knowit@mail.ru', 23, '6 second'::interval),
    ('kingdom', 'notnot123', 'king@mail.ru', 10, '7 second'::interval),
    ('lifeIsStrange', 'always', 'life@mail.ru', 111, '8 second'::interval),
    ('coala1', 'cheese', 'thesame@mail.ru', 12050, '9 second'::interval),
    ('prten1', 'apple', 'knowit@mail.ru', 23, '10 second'::interval),
    ('april', 'aprilapril', 'april@mail.ru', 10, '10 second'::interval),
    ('useruser', 'password', 'mail@mail.ru', 111, '10 second'::interval),
    ('tiger1', 'Bananas', 'tinan@mail.ru', 1000, '10 second'::interval),
    ('panda1', 'apple', 'today@mail.ru', 2323, '10 second'::interval),
    ('catmate1', 'juice', 'allday@mail.ru', 10000, '10 second'::interval),
    ('hotdog1', 'where', 'three@mail.ru', 88, '10 second'::interval),
    ('coala1', 'cheese', 'thesame@mail.ru', 12050, '10 second'::interval),
    ('prten1', 'apple', 'knowit@mail.ru', 23, '10 second'::interval),
    ('kingdom1', 'notnot123', 'king@mail.ru', 10, '10 second'::interval),
    ('lifeIsStrange1', 'always', 'life@mail.ru', 111, '10 second'::interval),
    ('coala11', 'cheese', 'thesame@mail.ru', 12050, '10 second'::interval),
    ('prten11', 'apple', 'knowit@mail.ru', 23, '10 second'::interval),
    ('april1', 'aprilapril', 'april@mail.ru', 10, '10 second'::interval),
    ('useruser1', 'password', 'mail@mail.ru', 111, '10 second'::interval),
    ('test1', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test2', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test3', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test4', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test5', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test6', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test7', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test8', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test9', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test10', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test11', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test12', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test13', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test14', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test15', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test16', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test17', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test18', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test19', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test20', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test21', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test22', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test23', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test24', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test25', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test26', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test27', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test28', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test29', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test30', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test40', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test41', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test42', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test43', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test44', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test45', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test46', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test47', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test48', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test49', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test50', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test51', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test52', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test53', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test54', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test55', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test56', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test57', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test58', 'password', 'mail@mail.ru', 0, '10 second'::interval),
    ('test59', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test60', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test61', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test62', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test63', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test64', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test65', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test66', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test67', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test68', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test69', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test70', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test71', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test72', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test73', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test74', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test75', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test76', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test77', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test78', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test79', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test80', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test81', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test82', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test83', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test84', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test85', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test86', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test87', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test88', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test89', 'password', 'mail@mail.ru', 0, '10 hour'::interval),
    ('test90', 'password', 'mail@mail.ru', 0, '10 hour'::interval);
*/

CREATE Table Record (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL,
    score int default 100,
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
    player_id   int NOT NULL,
    difficult int default 0,
    width   int NOT NULL,
    height   int NOT NULL,
    players   int NOT NULL,
    mines   int NOT NULL,
    date TIMESTAMPTZ not null,
    online bool
);

ALTER TABLE Game
ADD CONSTRAINT game_player
   FOREIGN KEY (player_id)
   REFERENCES Player(id)
   ON DELETE CASCADE;

CREATE Table Gamer (
    id SERIAL PRIMARY KEY,
    player_id int NOT NULL,
    game_id int NOT NULL,
    score int default 0,
    time interval default '24 hour'::interval,
    mines_open int default 0,
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

/*
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
    
        */
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
				Time:        ran.Intn(10000),
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
