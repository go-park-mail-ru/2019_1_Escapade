package main

import (
	"fmt"
	"net"
	"os"

	session "github.com/go-park-mail-ru/2019_1_Escapade/auth/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/router"

	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
)

func main() {
	// conf := config.Configure(confPath)
	// curlog := lg.Construct(logPath, logFile)
	// db := database.New(curlog)

	const (
		confPath = "auth/auth.json"
	)

	var (
		lis       net.Listener
		redisConn redis.Conn
		server    *grpc.Server
		conf      *config.Configuration
		err       error
	)

	if conf, err = config.InitPublic(confPath); err != nil {
		return
	}

	if lis, err = net.Listen("tcp", router.GetPort(conf)); err != nil {
		fmt.Println("cant listen that address:", err.Error())
		return
	}

	if redisConn, err = redis.DialURL(os.Getenv(conf.DataBase.URL)); err != nil {
		fmt.Println("cant connect to redis", err.Error())
		return
	}
	defer redisConn.Close()

	server = grpc.NewServer()
	session.RegisterAuthCheckerServer(server, session.NewSessionManager(redisConn, conf.Session))

	// curlog.Sugar.Infow("starting grpc server on "+conf.AC.Host+conf.AC.Port,
	// 	"source", "main.go")
	fmt.Println("Auth launched!")

	if err = server.Serve(lis); err != nil {
		fmt.Println("Auth catched error")
	}
}
