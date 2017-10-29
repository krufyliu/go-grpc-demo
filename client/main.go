package main

import (
	"log"

	"github.com/krufyliu/go-grpc-demo/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Authentication struct {
	Login    string
	Password string
}

// GetRequestMetadata gets the current request metadata
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"login":    a.Login,
		"password": a.Password,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security
func (a *Authentication) RequireTransportSecurity() bool {
	return false
}

func main() {
	auth := &Authentication{
		Login:    "liujun",
		Password: "1234",
	}
	conn, err := grpc.Dial("localhost:7777", grpc.WithInsecure(), grpc.WithPerRPCCredentials(auth))
	if err != nil {
		log.Fatalf("dial to localhost:7777 failed: %v", err)
	}
	defer conn.Close()
	rpcClient := api.NewPingClient(conn)
	req := api.PingMessage{Greeting: "hello"}
	rep, err := rpcClient.SayHello(context.Background(), &req)
	if err != nil {
		log.Fatalf("call SayHello failed: %v", err)
	}
	log.Printf("received message: %#v", rep.Greeting)
}
