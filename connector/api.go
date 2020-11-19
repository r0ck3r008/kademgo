package connector

import (
	"math/rand"
	"net"
	"time"

	"github.com/r0ck3r008/kademgo/utils"
)

// Ping requests for a ping reply from the passed in UDP address and returns a bool return value if
// a reply shows up. It waits for utils.PINGWAIT amount of time before it fails.
func (conn_p *Connector) Ping(srchash *[utils.HASHSZ]byte, addr_p *net.UDPAddr) bool {
	var rand_num int64 = int64(rand.Int())
	addr := *addr_p
	var pkt Pkt = Pkt{Ttl: utils.MAXHOPS, RandNum: rand_num, Hash: *srchash, Type: PingReq}
	var env Envelope = Envelope{pkt, addr}
	conn_p.sch <- env
	time.Sleep(utils.PINGWAIT)
	// Fetch result from map
	var ret bool = false
	conn_p.mut.Lock()
	if _, ok := conn_p.pcache[rand_num]; ok {
		ret = true
		delete(conn_p.pcache, rand_num)
	}
	conn_p.mut.Unlock()

	return ret
}

func (conn_p *Connector) FindPeers(srchash *[utils.HASHSZ]byte, gway_addr *string) {
	var gway_addr_p net.UDPAddr = net.UDPAddr{IP: []byte(*gway_addr), Port: utils.PORTNUM, Zone: ""}
	var rand_num int64 = int64(rand.Int())
	var pkt Pkt = Pkt{Ttl: utils.MAXHOPS, RandNum: rand_num, Hash: *srchash, Type: PingReq}
	var env Envelope = Envelope{pkt, gway_addr_p}
	conn_p.sch <- env
}
