package kademgo

import (
	"net"

	"github.com/r0ck3r008/kademgo/node"
)

type KademGo struct {
	node *node.Node
}

func (kdm_p *KadmeGo) Init(addr_p *net.UDPAddr) {
	var bind_addr string = "127.0.0.1"

	kdm_p.node = &Node{}
	kdm_p.node.Init(&bind_addr, addr_p)
}

func (kdm_p *KadmeGo) DeInit() {
	kdm_p.node.DeInit()
}
