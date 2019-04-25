package api

import (
	"escapade/internal/config"
	"escapade/internal/database"
	"fmt"
	"log"
	"os"
	"time"

	session "escapade/internal/services/auth/proto"

	"google.golang.org/grpc"
)

// Init creates Handler
func Init(DB *database.DataBase, c *config.Configuration) (handler *Handler) {

	if os.Getenv("AUTH_URL") == "" {
		os.Setenv("AUTH_URL", "localhost:3333")
	}
	grcpConn, err := grpc.Dial(
		os.Getenv("AUTH_URL"),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	//defer grcpConn.Close()

	sessManager := session.NewAuthCheckerClient(grcpConn)
	ws := config.WebSocketSettings{
		WriteWait:      time.Duration(c.WebSocket.WriteWait) * time.Second,
		PongWait:       time.Duration(c.WebSocket.PongWait) * time.Second,
		PingPeriod:     time.Duration(c.WebSocket.PingPeriod) * time.Second,
		MaxMessageSize: c.WebSocket.MaxMessageSize,
	}
	handler = &Handler{
		DB:              *DB,
		Storage:         c.Storage,
		Cookie:          c.Cookie,
		GameConfig:      c.Game,
		AWS:             c.AWS,
		WebSocket:       ws,
		WriteBufferSize: c.Server.WriteBufferSize,
		ReadBufferSize:  c.Server.ReadBufferSize,
		sessionManager:  sessManager,
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
