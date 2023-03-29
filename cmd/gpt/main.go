package main

import (
	"google.golang.org/grpc"
	"net"
	"net/http"
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
	go func() {
		halloPort := ":8000"
		log.Info("starting http hallo server ,port is ", halloPort)
		http.HandleFunc("/chat", server.ChatHandle)
		http.HandleFunc("/generateImage", server.GenerateImage)
		http.HandleFunc("/hallo", server.HalloHandler().ServeHTTP)
		http.ListenAndServe(halloPort, nil)
	}()

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
