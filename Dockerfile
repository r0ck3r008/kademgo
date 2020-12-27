FROM golang:alpine

# Install git
RUN apk add --no-cache git

# Clone Repo
RUN git clone https://github.com/r0ck3r008/kademgo

WORKDIR kademgo
