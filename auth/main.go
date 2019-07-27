package main

import (
	"net"
	"os"

	session "github.com/go-park-mail-ru/2019_1_Escapade/auth/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/server"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/utils"

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
		place     = "auth service:"
	)

	if conf, err = config.InitPublic(confPath); err != nil {
		return
	}

	if lis, err = net.Listen("tcp", router.GetPort(conf)); err != nil {
		utils.Debug(false, place, "cant listen address", err.Error())
		return
	}

	if redisConn, err = redis.DialURL(os.Getenv(conf.DataBase.URL)); err != nil {
		utils.Debug(false, place, "cant connect to redis", err.Error())
		return
	}
	defer redisConn.Close()

	server = grpc.NewServer()
	session.RegisterAuthCheckerServer(server, session.NewSessionManager(redisConn, conf.Session))

	utils.Debug(false, place, "success start!")

	if err = server.Serve(lis); err != nil {
		utils.Debug(false, place, "finish with error!", err.Error())
	}
}
