package main

import (
	"log"

	"github.com/krufyliu/go-grpc-demo/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:7777", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("dial to localhost:7777 failed: %v", err)
	}
	rpcClient := api.NewPingClient(conn)
	req := api.PingMessage{Greeting: "hello"}
	rep, err := rpcClient.SayHello(context.Background(), &req)
	if err != nil {
		log.Fatalf("call SayHello failed: %v", err)
	}
	log.Printf("received message: %#v", rep.Greeting)
}
