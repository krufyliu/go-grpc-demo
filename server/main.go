package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/krufyliu/go-grpc-demo/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// private type for Context keys
type contextKey int

const (
	clientIDKey contextKey = iota
)

func authenticateClient(ctx context.Context, s *api.Server) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing credentials")
	}
	clientLogin := strings.Join(md["login"], "")
	clientPassword := strings.Join(md["password"], "")
	if clientLogin != "liujun" {
		return "", fmt.Errorf("unknown user %s", clientLogin)
	}
	if clientPassword != "1234" {
		return "", fmt.Errorf("bad password %s", clientPassword)
	}
	log.Printf("authenticated client: %s", clientLogin)
	return "42", nil
}

func unaryInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	server, ok := info.Server.(*api.Server)
	if !ok {
		return nil, fmt.Errorf("unable to cast server")
	}
	clientID, err := authenticateClient(ctx, server)
	if err != nil {
		return nil, err
	}
	ctx = context.WithValue(ctx, clientIDKey, clientID)
	return handler(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Printf("failed to listne: %v", err)
	}
	log.Print("listen at :7777")
	s := api.Server{}
	rpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	api.RegisterPingServer(rpcServer, &s)
	err = rpcServer.Serve(lis)
	if err != nil {
		log.Printf("rpc serve failed: %v", err)
	}
}
