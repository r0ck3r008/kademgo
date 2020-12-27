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
func (conn_p *Connector) FindNodeReq(srchash, target *[utils.HASHSZ]byte, list *[]pkt.ObjAddr) pkt.ObjAddr {
	var ret pkt.ObjAddr
	var flag bool = false
	var oaddrs_list []*[]pkt.ObjAddr = []*[]pkt.ObjAddr{list}
	oaddr_blist := make(map[[utils.HASHSZ]byte]bool)
	for len(oaddrs_list) != 0 || !flag {
		var count int = 0
		var indx int = len(oaddrs_list) - 1
		var rands []int64
		for _, oaddr := range *(oaddrs_list[indx]) {
			if oaddr.Hash == *target {
				ret = oaddr
				flag = true
				break
			} else if count == utils.ALPHAVAL {
				break
			}
			if _, ok := oaddr_blist[oaddr.Hash]; !ok {
				rands = append(rands, int64(rand.Int()))
				addr := net.UDPAddr{IP: oaddr.Addr, Port: utils.PORTNUM, Zone: ""}
				var packet pkt.Packet = pkt.Packet{RandNum: rands[len(rands)-1], Hash: *srchash, Type: pkt.FindReq, THash: *target}
				var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: addr}
				conn_p.sch <- env

				oaddr_blist[oaddr.Hash] = true
				count++
			}
		}
		oaddrs_list = oaddrs_list[:indx-1]

		if !flag {
			time.Sleep(utils.PINGWAIT)
			for _, rand := range rands {
				conn_p.rwlock.RLock()
				if _, ok := conn_p.pcache[rand]; ok {
					var list_new [utils.KVAL]pkt.ObjAddr = conn_p.pcache[rand].Pkt.ObjArr
					var slice_new []pkt.ObjAddr = list_new[:]
					conn_p.rwlock.Lock()
					delete(conn_p.pcache, rand)
					conn_p.rwlock.Unlock()
					oaddrs_list = append(oaddrs_list, &slice_new)
				}
				conn_p.rwlock.RUnlock()
			}
		}
	}
	return ret
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
