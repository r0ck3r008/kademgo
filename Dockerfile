FROM golang:alpine

# Install git
RUN apk add --no-cache git make protoc

# Install grpc related modules
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go
RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc

# Clone Repo
RUN git clone https://github.com/r0ck3r008/kademgo
RUN make -C kademgo

WORKDIR kademgo
