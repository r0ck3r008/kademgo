package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/r0ck3r008/kademgo"
	"github.com/r0ck3r008/kademgo/protos"
	"github.com/r0ck3r008/kademgo/utils"
	"google.golang.org/grpc"
)

func main() {
	node := &kademgo.KademGo{}
	var gwayAddr string
	flag.StringVar(&gwayAddr, "gateway", "", "The IP address of gateway node")
	flag.Parse()

	if gwayAddr == "" {
		node.Init(nil, nil)
	} else {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", gwayAddr, utils.GRPCPORTNUM), grpc.WithInsecure())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in connecting to gateway grpc: %s\n", err)
			os.Exit(1)
		}
		gclient := protos.NewKademgoClient(conn)
		req := &protos.Request{}
		res, err := gclient.GetHash(context.Background(), req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in getting response: %s\n", err)
			os.Exit(1)
		}
		var hash [utils.HASHSZ]byte
		copy(hash[:], res.Hash[:])
		node.Init(&gwayAddr, &hash)
	}

	fmt.Printf("Sleeping before DeInit")
	time.Sleep(1)
	node.DeInit()
}
