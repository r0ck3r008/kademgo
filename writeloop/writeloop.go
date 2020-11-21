package writeloop

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/r0ck3r008/kademgo/objmap"
	"github.com/r0ck3r008/kademgo/utils"
)

type WriteLoop struct {
	sch    chan utils.Envelope
	pcache *map[int64]utils.Envelope
	mut    *sync.Mutex
}

func (wrl_p *WriteLoop) Init(mut *sync.Mutex, pcache *map[int64]utils.Envelope) {
	wrl_p.sch = make(chan utils.Envelope, 100)
	wrl_p.mut = mut
	wrl_p.pcache = pcache
}

func (wrl_p *WriteLoop) DeInit() {
	close(wrl_p.sch)
}

// WriteLoop is supposed to be run as a goroutine which takes all the packets that need to be sent
// from the node and send them asynchronously to the desired destinations.
func (wrl_p *WriteLoop) WriteLoop(conn_p *net.UDPConn) {
	for env := range wrl_p.sch {
		cmds, err := json.Marshal(env.pkt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Marshalling: %s\n", err)
			break
		}
		if _, err := conn_p.WriteToUDP(cmds, &env.addr); err != nil {
			fmt.Fprintf(os.Stderr, "Error in writing: %s\n", err)
			break
		}
	}
}

// Ping requests for a ping reply from the passed in UDP address and returns a bool return value if
// a reply shows up. It waits for utils.PINGWAIT amount of time before it fails.
func (wrl_p *WriteLoop) Ping(srchash *[utils.HASHSZ]byte, addr_p *net.UDPAddr) bool {
	var rand_num int64 = int64(rand.Int())
	addr := *addr_p
	var pkt utils.Pkt = utils.Pkt{Ttl: utils.MAXHOPS, RandNum: rand_num, Hash: *srchash, Type: PingReq}
	var env utils.Envelope = utils.Envelope{pkt, addr}
	wrl_p.sch <- env
	time.Sleep(utils.PINGWAIT)
	// Fetch result from map
	var ret bool = false
	wrl_p.mut.Lock()
	if _, ok := wrl_p.pcache[rand_num]; ok {
		ret = true
		delete(wrl_p.pcache, rand_num)
	}
	wrl_p.mut.Unlock()

	return ret
}

// Store is a shoot and forget type of a function. It works in best effort way.
func (wrl_p *WriteLoop) Store(srchash *[utils.HASHSZ]byte, addr_p *net.UDPAddr, obj_p *objmap.ObjNode) {
	var rand_num int64 = int64(rand.Int())
	obj := *obj_p
	var pkt utils.Pkt = utils.Pkt{Ttl: utils.MAXHOPS, RandNum: rand_num, Type: Store, Hash: *srchash, Obj: obj}
	var env utils.Envelope = utils.Envelope{pkt, *addr_p}

	wrl_p.sch <- env
}
