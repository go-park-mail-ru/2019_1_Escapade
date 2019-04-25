package main

import (
	sessMan "escapade/internal/services/auth"
	session "escapade/internal/services/auth/proto"
	"fmt"
	"log"
	"net"

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

	redisConn, err := redis.DialURL("redis://user:@localhost:6379/0")
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
