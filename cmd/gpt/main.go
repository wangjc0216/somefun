package main

import (
	"google.golang.org/grpc"
	"net"
	"somefun/api/gpt"
	"somefun/pkg/conf"
	"somefun/pkg/log"
	"somefun/pkg/server"
)

func init() {
	log.Init("gpt-api.log", "info", "gpt-api", true)
}

func main() {

	conf.LoadConfig("./config.json")

	gptHandler, err := server.NewGptHandler()
	if err != nil {
		log.Fatalf("new gpt handler error:", err)
	}
	s := grpc.NewServer()
	gpt.RegisterGptServiceServer(s, gptHandler)
	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("listener error:", err)
	}
	log.Info("starting grpc server ,port is ", port)
	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("server error:", err)
	}

}
