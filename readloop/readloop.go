package readloop

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	sync "sync"

	"github.com/r0ck3r008/kademgo/nbrmap"
	"github.com/r0ck3r008/kademgo/objmap"
	"github.com/r0ck3r008/kademgo/utils"
)

type ReadLoop struct {
	rch    chan utils.Envelope
	sch    chan<- utils.Envelope
	pcache *map[int64]utils.Envelope
	mut    *sync.Mutex
	nmap   *nbrmap.NbrMap
	ost    *objmap.ObjMap
}

func (rdl_p *ReadLoop) Init(mut *sync.Mutex, pcache *map[int64]utils.Envelope, sch chan<- utils.Envelope) {
	rdl_p.rch = make(chan Envelope, 100)
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
		switch env.pkt.Type {
		case PingReq:
			wg.Add(1)
			go func() { rdl_p.PingRes(env); wg.Done() }()
		case Store:
			wg.Add(1)
			go func() { rdl_p.StoreHandler(env); wg.Done() }()
		default:
			// Acquire write lock and write to cache
			rdl_p.mut.Lock()
			rdl_p.pcache[env.pkt.RandNum] = env
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
		var pkt Pkt = Pkt{}
		err = json.Unmarshal(cmdr, &pkt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in unmarshalling: %s\n", err)
			os.Exit(1)
		}
		if pkt.Ttl != 0 {
			pkt.Ttl--
			var env utils.Envelope = utils.Envelope{pkt, *addr_p}
			rdl_p.rch <- env
		}
	}
}

func (rdl_p *ReadLoop) PingRes(env utils.Envelope) {
	env.pkt.Type = PingRes
	rdl_p.sch <- env

	// Insert to the NbrMap here.
}

func (rdl_p *ReadLoop) StoreHandler(env utils.Envelope) {
	// Insert to the ObjMap here
}
