package main

import (
	sessMan "escapade/internal/services/auth"
	session "escapade/internal/services/auth/proto"
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

	lis, err := net.Listen("tcp", ":3334")
	if err != nil {
		fmt.Println(err)
	}

	server := grpc.NewServer()
	session.RegisterAuthCheckerServer(server, sessMan.NewSessionManager(redisConn))

	// curlog.Sugar.Infow("starting grpc server on "+conf.AC.Host+conf.AC.Port,
	// 	"source", "main.go")
	fmt.Println("Scuccess")

	server.Serve(lis)
}
