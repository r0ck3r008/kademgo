package connector

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
)

// writeloop is supposed to be run as a goroutine which takes all the packets that need to be sent
// from the node and send them asynchronously to the desired destinations.
func (conn_p *Connector) writeloop() {
	for env := range conn_p.sch {
		cmds, err := json.Marshal(env.Pkt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Marshalling: %s\n", err)
			break
		}
		if _, err := conn_p.conn.WriteToUDP(cmds, &env.Addr); err != nil {
			fmt.Fprintf(os.Stderr, "Error in writing: %s\n", err)
			break
		}
	}
}

// PingReq requests for a ping reply from the passed in UDP address and returns a bool return value if
// a reply shows up. It waits for utils.PINGWAIT amount of time before it fails.
// It is called by NbrMap while deciding whether to evict a least seen nbr from cache
func (conn_p *Connector) PingReq(srchash *[utils.HASHSZ]byte, addr_p *net.IP) bool {
	var rand_num int64 = int64(rand.Int())
	addr := net.UDPAddr{IP: *addr_p, Port: utils.PORTNUM, Zone: ""}
	var packet pkt.Packet = pkt.Packet{RandNum: rand_num, Hash: *srchash, Type: pkt.PingReq}
	var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: addr}
	conn_p.sch <- env
	time.Sleep(utils.PINGWAIT)
	// Fetch result from map
	var ret bool = false
	conn_p.rwlock.RLock()
	if _, ok := conn_p.pcache[rand_num]; ok {
		ret = true
		conn_p.rwlock.Lock()
		delete(conn_p.pcache, rand_num)
		conn_p.rwlock.Unlock()
	}
	conn_p.rwlock.RUnlock()

	return ret
}

// PingRes is responsible for writing the Ping Response to the network.
// This is called by the Node.PingReqHandler.
func (conn_p *Connector) PingRes(env *pkt.Envelope) {
	conn_p.sch <- *env
}

// FindNode is called by Node and expects ALPHAVAL number of nbrs to recursively
// query for neighbours.
func (conn_p *Connector) FindNodeReq(srchash, target *[utils.HASHSZ]byte, ret *[]pkt.ObjAddr) {
	var rands []int64
	for i := 0; i < utils.ALPHAVAL; i++ {
		rands = append(rands, int64(rand.Int()))
		addr := net.UDPAddr{IP: (*ret)[i].Addr, Port: utils.PORTNUM, Zone: ""}
		var packet pkt.Packet = pkt.Packet{RandNum: rands[i], Hash: *srchash, Type: pkt.FindReq, THash: *target}
		var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: addr}
		conn_p.sch <- env
	}
	time.Sleep(utils.PINGWAIT)
	// Continue RPC with received values
}

// FindNodeRes is responsible for writing the KVAL Nbrs found by Node.FindReqHandler to the network.
func (conn_p Connector) FindNodeRes(srchash *[utils.HASHSZ]byte, env *pkt.Envelope, ret *[]pkt.ObjAddr) {
	var objarr [utils.KVAL]pkt.ObjAddr
	copy(objarr[:], (*ret)[:])
	var pkt pkt.Packet = pkt.Packet{RandNum: (*env).Pkt.RandNum, Hash: *srchash, Type: pkt.FindRes, ObjArr: objarr}

	// Send out the reply
	(*env).Pkt = pkt
	conn_p.sch <- *env
}
