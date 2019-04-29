package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_Escapade/internal/config"
	sessMan "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth"
	session "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/proto"
	"net"
	"os"

	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
)

func main() {
	// conf := config.Configure(confPath)
	// curlog := lg.Construct(logPath, logFile)
	// db := database.New(curlog)

	var (
		lis       net.Listener
		redisConn redis.Conn
		server    *grpc.Server
		conf      *config.Configuration
		err       error
	)

	if lis, err = net.Listen("tcp", ":3333"); err != nil {
		fmt.Println(err)
	}
	if os.Getenv("REDIS_URL") == "" {
		os.Setenv("REDIS_URL", "redis://user:@localhost:6379/0")
	}

	if redisConn, err = redis.DialURL(os.Getenv("REDIS_URL")); err != nil {
		fmt.Println("cant connect to redis")
		return
	}
	defer redisConn.Close()
	const (
		confPath = "conf.json"
	)
	if conf, err = config.Init(confPath, ""); err != nil {
		return
	}

	server = grpc.NewServer()
	session.RegisterAuthCheckerServer(server, sessMan.NewSessionManager(redisConn, conf.Session))

	// curlog.Sugar.Infow("starting grpc server on "+conf.AC.Host+conf.AC.Port,
	// 	"source", "main.go")
	fmt.Println("Auth launched!")

	if err = server.Serve(lis); err != nil {
		fmt.Println("Auth catched error")
	}
}
