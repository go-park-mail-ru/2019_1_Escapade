package api

import (
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/database"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/models"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

	"context"
	ran "math/rand"
	
	session "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/proto"


	"fmt"
	"time"
)

// Init creates Handler
func Init(DB *database.DataBase, c *config.Configuration) (handler *Handler) {
	ws := config.WebSocketSettings{
		WriteWait:      time.Duration(c.WebSocket.WriteWait) * time.Second,
		PongWait:       time.Duration(c.WebSocket.PongWait) * time.Second,
		PingPeriod:     time.Duration(c.WebSocket.PingPeriod) * time.Second,
		MaxMessageSize: c.WebSocket.MaxMessageSize,
	}
	handler = &Handler{
		DB:              *DB,
		Storage:         c.Storage,
		Session:         c.Session,
		GameConfig:      c.Game,
		AWS:             c.AWS,
		WebSocket:       ws,
		WriteBufferSize: c.Server.WriteBufferSize,
		ReadBufferSize:  c.Server.ReadBufferSize,
	}
	return
}

// GetHandler return created handler with database and configuration
func GetHandler(confPath, secretPath string) (handler *Handler,
	conf *config.Configuration, err error) {

	var (
		db *database.DataBase
	)

	if conf, err = config.Init(confPath, secretPath); err != nil {
		return
	}
	fmt.Println("confPath done")

	if db, err = database.Init(conf.DataBase); err != nil {
		return
	}

	fmt.Println("database done")
	handler = Init(db, conf)
	return
}

func (h *Handler) RandomUsers(limit int) {

	n := 16
	for i := 0; i < limit; i++ {
		ran.Seed(time.Now().UnixNano())
		user := &models.UserPrivateInfo{
			Name:     utils.RandomString(n),
			Email:    utils.RandomString(n),
			Password: utils.RandomString(n)}
		id, _ := h.DB.Register(user)
		ctx := context.Background()
		sessID, err := h.Clients.Session.Create(ctx, &session.Session{
			UserID: int32(id),
			Login:  user.Name,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("sessID:", sessID)

		for j := 0; j < 4; j++ {
			record := &models.Record{
				Score:       ran.Intn(1000000),
				Time:        float64(ran.Intn(10000)),
				Difficult:   j,
				SingleTotal: ran.Intn(2),
				OnlineTotal: ran.Intn(2),
				SingleWin:   ran.Intn(2),
				OnlineWin:   ran.Intn(2)}
			h.DB.UpdateRecords(id, record)
		}

	}
}
