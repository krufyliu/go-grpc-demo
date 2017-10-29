package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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

func credMatcher(headerName string) (mdName string, ok bool) {
	if headerName == "Login" || headerName == "Password" {
		return headerName, true
	}
	return "", false
}

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

func startGPRCServer(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("[grpc] failed to listen: %v", err)
	}
	log.Printf("[grpc] listen at %s", addr)
	s := api.Server{}
	rpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	api.RegisterPingServer(rpcServer, &s)
	err = rpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("[grpc] rpc serve failed: %v", err)
	}
	return nil
}

func startRestServer(addr, grpcAddress string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(credMatcher))
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	err := api.RegisterPingHandlerFromEndpoint(ctx, mux, grpcAddress, dialOpts)
	if err != nil {
		return fmt.Errorf("[rest] could not register service Ping: %s", err)
	}
	log.Printf("[rest] starting HTTP/1.1 REST server on %s", addr)
	return http.ListenAndServe(addr, mux)
}

func main() {
	grpcAddr := "localhost:7777"
	restAddr := "localhost:7778"
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		err := startGPRCServer(grpcAddr)
		if err != nil {
			log.Printf("failed to start gRPC server: %s", err)
		}
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		err := startRestServer(restAddr, grpcAddr)
		if err != nil {
			log.Printf("failed to start rest server: %s", err)
		}
	}()
	log.Print("wait...")
	wg.Wait()
}
