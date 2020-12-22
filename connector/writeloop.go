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

// Ping requests for a ping reply from the passed in UDP address and returns a bool return value if
// a reply shows up. It waits for utils.PINGWAIT amount of time before it fails.
func (conn_p *Connector) Ping(srchash *[utils.HASHSZ]byte, addr_p *net.IP) bool {
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

// Store is a shoot and forget type of a function. It works in best effort way.
func (conn_p *Connector) Store(srchash *[utils.HASHSZ]byte, addr_p *net.IP, obj_p *pkt.ObjAddr) {
	var rand_num int64 = int64(rand.Int())
	obj := *obj_p
	addr := net.UDPAddr{IP: *addr_p, Port: utils.PORTNUM, Zone: ""}
	var packet pkt.Packet = pkt.Packet{RandNum: rand_num, Type: pkt.Store, Hash: *srchash, Obj: obj}
	var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: addr}

	conn_p.sch <- env
}
