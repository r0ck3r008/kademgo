package main

import (
	"time"

	"github.com/r0ck3r008/kademgo"
)

func main() {
	var node *kademgo.KademGo = &kademgo.KademGo{}
	node.Init(nil, nil)

	// Give time for Readloop to begin before sending deInit
	time.Sleep(time.Second)
	node.DeInit()
}
