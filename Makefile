all: go

go: protos
	go build

.PHONY: protos
protos:
	protoc -I protos/ protos/kademgo.proto --go_out=protos
	protoc -I protos/ protos/kademgo.proto --go-grpc_out=require_unimplemented_servers=false:protos

clean:
	rm -f protos/*.pb.go