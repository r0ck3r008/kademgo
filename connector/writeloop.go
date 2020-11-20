package connector

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/r0ck3r008/kademgo/objmap"
	"github.com/r0ck3r008/kademgo/utils"
)

// WriteLoop is supposed to be run as a goroutine which takes all the packets that need to be sent
// from the node and send them asynchronously to the desired destinations.
func (conn_p *Connector) WriteLoop() {
	for env := range conn_p.sch {
		cmds, err := json.Marshal(env.pkt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Marshalling: %s\n", err)
			break
		}
		if _, err := conn_p.conn.WriteToUDP(cmds, &env.addr); err != nil {
			fmt.Fprintf(os.Stderr, "Error in writing: %s\n", err)
			break
		}
	}
}

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

// Store is a shoot and forget type of a function. It works in best effort way.
func (conn_p *Connector) Store(srchash *[utils.HASHSZ]byte, addr_p *net.UDPAddr, obj_p *objmap.ObjNode) {
	var rand_num int64 = int64(rand.Int())
	obj := *obj_p
	var pkt Pkt = Pkt{Ttl: utils.MAXHOPS, RandNum: rand_num, Type: Store, Hash: *srchash, Obj: obj}
	var env Envelope = Envelope{pkt, *addr_p}

	conn_p.sch <- env
}
