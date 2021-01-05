.PHONY: protos

protos:
	protoc -I protos/ protos/kademgo.proto --go_out=protos
	protoc -I protos/ protos/kademgo.proto --go-grpc_out=protos

clean:
	rm -f protos/*.pb.go