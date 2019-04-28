package main

import (
	sessMan "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth"
	session "github.com/go-park-mail-ru/2019_1_Escapade/internal/services/auth/proto"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
)

func main() {
	// conf := config.Configure(confPath)
	// curlog := lg.Construct(logPath, logFile)
	// db := database.New(curlog)

	lis, err := net.Listen("tcp", ":3333")
	if err != nil {
		fmt.Println(err)
	}
	if os.Getenv("REDIS_URL") == "" {
		os.Setenv("REDIS_URL", "redis://user:@localhost:6379/0")
	}
	redisConn, err := redis.DialURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("cant connect to redis")
	}
	server := grpc.NewServer()
	session.RegisterAuthCheckerServer(server, sessMan.NewSessionManager(redisConn))

	// curlog.Sugar.Infow("starting grpc server on "+conf.AC.Host+conf.AC.Port,
	// 	"source", "main.go")
	fmt.Println("Scuccess")

	server.Serve(lis)
}
