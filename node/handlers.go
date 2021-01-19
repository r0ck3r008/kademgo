package node

import (
	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
)

// PingReqHandler updates the packet to indicate it is a ping response
// and calls the Connector.PingRes for writing to network. It also updates
// the NbrMap cache for the received Nbr.
func (nodeP *Node) PingReqHandler(env pkt.Envelope) {
	// Send response
	env.Pkt.Type = pkt.PingRes
	nodeP.conn.PingRes(&env)

	// Potentially update in NbrMap
	nodeP.nmap.Insert(&nodeP.hash, &env.Pkt.Hash, &env.Addr.IP, nodeP.conn)
}

// FindReqHandler is responsible for finding KVAL Nbrs from NbrMap and calling
// Connector.FindNodeRes for writing the response on the network.
func (nodeP *Node) FindReqHandler(env pkt.Envelope) {
	var ret []pkt.ObjAddr
	nodeP.nmap.NodeLookup(&nodeP.hash, &env.Pkt.THash, &ret, utils.KVAL)

	nodeP.conn.FindNodeRes(&nodeP.hash, &env, &ret)
}
