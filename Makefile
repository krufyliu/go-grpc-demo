.PHONY: all api clean dep help client server

all: server client

proto: api/*.proto
	protoc -I api --go_out=plugins=grpc:api api/*.proto

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