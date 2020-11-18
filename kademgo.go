package kademgo

import (
	"fmt"
	"os"

	"github.com/r0ck3r008/kademgo/node"
)

// KademInit is the function that initiates the KademGo Library
func KademInit() {
	var bind_addr string = "127.0.0.1"
	_, err := node.NodeInit(&bind_addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	fmt.Println("Kademgo")
}
