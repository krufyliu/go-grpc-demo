package main

import (
	"log"
	"net"

	"github.com/krufyliu/go-grpc-demo/api"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Printf("failed to listne: %v", err)
	}
	log.Print("listen at :7777")
	s := api.Server{}
	rpcServer := grpc.NewServer()
	api.RegisterPingServer(rpcServer, &s)
	err = rpcServer.Serve(lis)
	if err != nil {
		log.Printf("rpc serve failed: %v", err)
	}
}
