package main

import (
	"github.com/r0ck3r008/kademgo"
)

func main() {
	var node *kademgo.KademGo = &kademgo.KademGo{}
	node.Init(nil, nil)
	defer node.DeInit()
}
