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
func (connP *Connector) writeloop() {
	for env := range connP.sch {
		cmds, err := json.Marshal(env.Pkt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Marshalling: %s\n", err)
			break
		}
		if _, err := connP.conn.WriteToUDP(cmds, &env.Addr); err != nil {
			fmt.Fprintf(os.Stderr, "Error in writing: %s\n", err)
			break
		}
	}
}

// PingReq requests for a ping reply from the passed in UDP address and returns a bool return value if
// a reply shows up. It waits for utils.PINGWAIT amount of time before it fails.
// It is called by NbrMap while deciding whether to evict a least seen nbr from cache
func (connP *Connector) PingReq(srchash *[utils.HASHSZ]byte, addrP *net.IP) bool {
	var randNum int64 = int64(rand.Int())
	addr := net.UDPAddr{IP: *addrP, Port: utils.PORTNUM, Zone: ""}
	var packet pkt.Packet = pkt.Packet{RandNum: randNum, Hash: *srchash, Type: pkt.PingReq}
	var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: addr}
	connP.sch <- env
	time.Sleep(utils.TIMEOUT)
	// Fetch result from map
	var ret bool = false
	connP.rwlock.RLock()
	if _, ok := connP.pcache[randNum]; ok {
		ret = true
		connP.rwlock.Lock()
		delete(connP.pcache, randNum)
		connP.rwlock.Unlock()
	}
	connP.rwlock.RUnlock()

	return ret
}

// PingRes is responsible for writing the Ping Response to the network.
// This is called by the Node.PingReqHandler.
func (connP *Connector) PingRes(env *pkt.Envelope) {
	connP.sch <- *env
}

// FindNodeReq is called by Node and expects ALPHAVAL number of nbrs to recursively
// query for neighbours.
func (connP *Connector) FindNodeReq(srchash, target *[utils.HASHSZ]byte, list *[]pkt.ObjAddr) pkt.ObjAddr {
	var ret pkt.ObjAddr
	var flag bool = false
	var oaddrsList []*[]pkt.ObjAddr = []*[]pkt.ObjAddr{list}
	oaddrBlist := make(map[[utils.HASHSZ]byte]bool)
	for len(oaddrsList) != 0 || !flag {
		var count int = 0
		var indx int = len(oaddrsList) - 1
		var rands []int64
		for _, oaddr := range *(oaddrsList[indx]) {
			if oaddr.Hash == *target {
				ret = oaddr
				flag = true
				break
			} else if count == utils.ALPHAVAL {
				break
			}
			if _, ok := oaddrBlist[oaddr.Hash]; !ok {
				rands = append(rands, int64(rand.Int()))
				addr := net.UDPAddr{IP: oaddr.Addr, Port: utils.PORTNUM, Zone: ""}
				var packet pkt.Packet = pkt.Packet{RandNum: rands[len(rands)-1], Hash: *srchash, Type: pkt.FindReq, THash: *target}
				var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: addr}
				connP.sch <- env

				oaddrBlist[oaddr.Hash] = true
				count++
			}
		}
		oaddrsList = oaddrsList[:indx-1]

		if !flag {
			time.Sleep(utils.TIMEOUT)
			for _, rand := range rands {
				connP.rwlock.RLock()
				if _, ok := connP.pcache[rand]; ok {
					var listNew [utils.KVAL]pkt.ObjAddr = connP.pcache[rand].Pkt.ObjArr
					var sliceNew []pkt.ObjAddr = listNew[:]
					connP.rwlock.Lock()
					delete(connP.pcache, rand)
					connP.rwlock.Unlock()
					oaddrsList = append(oaddrsList, &sliceNew)
				}
				connP.rwlock.RUnlock()
			}
		}
	}
	return ret
}

// FindNodeRes is responsible for writing the KVAL Nbrs found by Node.FindReqHandler to the network.
func (connP Connector) FindNodeRes(srchash *[utils.HASHSZ]byte, env *pkt.Envelope, ret *[]pkt.ObjAddr) {
	var objarr [utils.KVAL]pkt.ObjAddr
	copy(objarr[:], (*ret)[:])
	var pkt pkt.Packet = pkt.Packet{RandNum: (*env).Pkt.RandNum, Hash: *srchash, Type: pkt.FindRes, ObjArr: objarr}

	// Send out the reply
	(*env).Pkt = pkt
	connP.sch <- *env
}
