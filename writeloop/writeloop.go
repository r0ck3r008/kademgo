// writeloop is the loop that gets the packets that need to be send out
// per requestors all over the library.
package writeloop

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
)

// WriteLoop is the handle that each sender that needs to put out anything to the wire
// gets and eventually calls one of the methods of WriteLoop.
type WriteLoop struct {
	pcache *map[int64]pkt.Envelope
	mut    *sync.RWMutex
	// sch has only ony sink and a lot of sources.
	// The sink is WriteLoop
	sch chan pkt.Envelope
}

// Init initiates all the internal members of WriteLoop.
func (wrl_p *WriteLoop) Init(mut *sync.RWMutex, pcache *map[int64]pkt.Envelope, sch chan pkt.Envelope) {
	wrl_p.mut = mut
	wrl_p.pcache = pcache
	wrl_p.sch = sch
}

// WriteLoop is supposed to be run as a goroutine which takes all the packets that need to be sent
// from the node and send them asynchronously to the desired destinations.
func (wrl_p *WriteLoop) WriteLoop(conn_p *net.UDPConn) {
	for env := range wrl_p.sch {
		cmds, err := json.Marshal(env.Pkt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Marshalling: %s\n", err)
			break
		}
		if _, err := conn_p.WriteToUDP(cmds, &env.Addr); err != nil {
			fmt.Fprintf(os.Stderr, "Error in writing: %s\n", err)
			break
		}
	}
}

// Ping requests for a ping reply from the passed in UDP address and returns a bool return value if
// a reply shows up. It waits for utils.PINGWAIT amount of time before it fails.
func (wrl_p *WriteLoop) Ping(srchash *[utils.HASHSZ]byte, addr_p *net.IP) bool {
	var rand_num int64 = int64(rand.Int())
	addr := net.UDPAddr{IP: *addr_p, Port: utils.PORTNUM, Zone: ""}
	var packet pkt.Packet = pkt.Packet{RandNum: rand_num, Hash: *srchash, Type: pkt.PingReq}
	var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: addr}
	wrl_p.sch <- env
	time.Sleep(utils.PINGWAIT)
	// Fetch result from map
	var ret bool = false
	wrl_p.mut.RLock()
	if _, ok := (*wrl_p.pcache)[rand_num]; ok {
		ret = true
		wrl_p.mut.Lock()
		delete((*wrl_p.pcache), rand_num)
		wrl_p.mut.Unlock()
	}
	wrl_p.mut.RUnlock()

	return ret
}

// Store is a shoot and forget type of a function. It works in best effort way.
func (wrl_p *WriteLoop) Store(srchash *[utils.HASHSZ]byte, addr_p *net.IP, obj_p *pkt.ObjAddr) {
	var rand_num int64 = int64(rand.Int())
	obj := *obj_p
	addr := net.UDPAddr{IP: *addr_p, Port: utils.PORTNUM, Zone: ""}
	var packet pkt.Packet = pkt.Packet{RandNum: rand_num, Type: pkt.Store, Hash: *srchash, Obj: obj}
	var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: addr}

	wrl_p.sch <- env
}
