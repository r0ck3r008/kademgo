// readloop package implements the handler logic of the packets received by the peer.
package readloop

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	sync "sync"

	"github.com/r0ck3r008/kademgo/nbrmap"
	"github.com/r0ck3r008/kademgo/objmap"
	"github.com/r0ck3r008/kademgo/pkt"
)

// ReadLoop structure is the handler that node has and uses to asynchronously
// receive packets and process or cache.
type ReadLoop struct {
	rch    chan pkt.Envelope
	sch    chan<- pkt.Envelope
	pcache *map[int64]pkt.Envelope
	mut    *sync.RWMutex
	nmap   *nbrmap.NbrMap
	ost    *objmap.ObjMap
}

// Init initiates the internal members of the ReadLoop type.
func (rdl_p *ReadLoop) Init(mut *sync.RWMutex, pcache *map[int64]pkt.Envelope, sch chan<- pkt.Envelope) {
	rdl_p.rch = make(chan pkt.Envelope, 100)
	rdl_p.sch = sch
	rdl_p.mut = mut
	rdl_p.pcache = pcache

	rdl_p.nmap = &nbrmap.NbrMap{}
	rdl_p.ost = &objmap.ObjMap{}
	rdl_p.nmap.Init()
	rdl_p.ost.Init()
}

// DeInit closes the receive channel and makes the Collector exit.
func (rdl_p *ReadLoop) DeInit() {
	close(rdl_p.rch)
}

// Collector is intended to be a goroutine that process the received packets in form of Envelope
// struct and caches it in the connector cache based on the identifier.
func (rdl_p *ReadLoop) Collector() {
	wg := sync.WaitGroup{}
	for env := range rdl_p.rch {
		switch env.Pkt.Type {
		case pkt.PingReq:
			wg.Add(1)
			go func() { rdl_p.PingRes(env); wg.Done() }()
		case pkt.Store:
			wg.Add(1)
			go func() { rdl_p.StoreHandler(env); wg.Done() }()
		default:
			// Acquire write lock and write to cache
			rdl_p.mut.Lock()
			(*rdl_p.pcache)[env.Pkt.RandNum] = env
			rdl_p.mut.Unlock()
		}
	}
	wg.Wait()
}

// ReadLoop is supposed to be run as a go routine which can read all the messages comming in
// to the node and send those along, if the TTL has not expired, to the Collector.
func (rdl_p *ReadLoop) ReadLoop(conn_p *net.UDPConn) {
	for {
		var cmdr []byte
		_, addr_p, err := conn_p.ReadFromUDP(cmdr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in reading: %s\n", err)
			break
		}
		var packet pkt.Packet = pkt.Packet{}
		err = json.Unmarshal(cmdr, &packet)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in unmarshalling: %s\n", err)
			os.Exit(1)
		}
		var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: *addr_p}
		// BUG: This might make the application panic if DeInit on ReadLoop is called while
		// receive channel is being written to with a new packet.
		rdl_p.rch <- env
	}
}

// PingRes is the handler of the Ping Request that node receives.
func (rdl_p *ReadLoop) PingRes(env pkt.Envelope) {
	env.Pkt.Type = pkt.PingRes
	rdl_p.sch <- env

	// Insert to the NbrMap here.
}

func (rdl_p *ReadLoop) StoreHandler(env pkt.Envelope) {
	// Insert to the ObjMap here
}
