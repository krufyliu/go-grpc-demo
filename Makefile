.PHONY: all api clean dep help client server

all: server client

api/api.pb.go: api/api.proto
	protoc -I api/ \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:api \
	api/api.proto

api/api.pb.gw.go: api/api.proto
	protoc -I api/ \
    -I${GOPATH}/src \
    -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --grpc-gateway_out=logtostderr=true:api \
    api/api.proto
api/api.swagger.json:
	protoc -I api/ \
  	-I${GOPATH}/src \
  	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  	--swagger_out=logtostderr=true:api \
  	api/api.proto

api: api/api.pb.go api/api.pb.gw.go api/api.swagger.json

client: client/main.go
	go build -o build/client client/main.go

server: server/main.go
	go build -o build/server server/main.go

dep:
	go get -v -d ./...

clean:
	@rm -rf build

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'