CREATE TABLE Player (
    id SERIAL PRIMARY KEY,
    name varchar(30) NOT NULL,
    password varchar(30) NOT NULL,
    email varchar(30) NOT NULL,
    photo_title varchar(50),
    --FirstSeen   timestamp without time zone NOT NULL,
	--LastSeen    timestamp without time zone NOT NULL,
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