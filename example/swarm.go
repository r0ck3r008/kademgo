package main

import (
	"flag"

	"github.com/r0ck3r008/kademgo"
)

func makePeers(ret []kademgo.KademGo, npeers int) {
	for i := 0; i < npeers; i++ {
		if i == 0 {
			ret[i].Init(nil, nil)
		} else {
			// TODO: Fetch hash and IP address of (i-1)th peer and pass here
			ret[i].Init(nil, nil)
		}
	}
}

func main() {
	var npeers int
	flag.IntVar(&npeers, "npeers", 0, "The number of nodes in the network")
	flag.Parse()
	var nodes []kademgo.KademGo

	if npeers > 0 {
		makePeers(nodes, npeers)
	}
}
