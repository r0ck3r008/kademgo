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

type ReadLoop struct {
	rch    chan pkt.Envelope
	sch    chan<- pkt.Envelope
	pcache *map[int64]pkt.Envelope
	mut    *sync.Mutex
	nmap   *nbrmap.NbrMap
	ost    *objmap.ObjMap
}

func (rdl_p *ReadLoop) Init(mut *sync.Mutex, pcache *map[int64]pkt.Envelope, sch chan<- pkt.Envelope) {
	rdl_p.rch = make(chan pkt.Envelope, 100)
	rdl_p.sch = sch
	rdl_p.mut = mut
	rdl_p.pcache = pcache

	rdl_p.nmap = &nbrmap.NbrMap{}
	rdl_p.ost = &objmap.ObjMap{}
	rdl_p.nmap.Init()
	rdl_p.ost.Init()
}

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
		if packet.Ttl != 0 {
			packet.Ttl--
			var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: *addr_p}
			rdl_p.rch <- env
		}
	}
}

func (rdl_p *ReadLoop) PingRes(env pkt.Envelope) {
	env.Pkt.Type = pkt.PingRes
	rdl_p.sch <- env

	// Insert to the NbrMap here.
}

func (rdl_p *ReadLoop) StoreHandler(env pkt.Envelope) {
	// Insert to the ObjMap here
}
