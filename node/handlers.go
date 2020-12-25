package node

import (
	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
)

func (node_p *Node) PingReqHandler(env pkt.Envelope) {
	// Send response
	env.Pkt.Type = pkt.PingRes
	node_p.conn.PingRes(&env)

	// Potentially update in NbrMap
	node_p.nmap.Insert(node_p.hash, env.Pkt.Hash, env.Addr.IP, node_p.conn)
}

func (node_p *Node) FindReqHandler(env pkt.Envelope) {
	var ret []pkt.ObjAddr
	node_p.nmap.NodeLookup(&node_p.hash, &env.Pkt.THash, &ret, utils.KVAL)

	node_p.conn.FindNodeRes(&node_p.hash, &env, &ret)
}
