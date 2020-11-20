package connector

import (
	"encoding/json"
	"fmt"
	"os"
	sync "sync"
)

// Collector is intended to be a goroutine that process the received packets in form of Envelope
// struct and caches it in the connector cache based on the identifier.
func (conn_p *Connector) Collector() {
	wg := sync.WaitGroup{}
	for env := range conn_p.rch {
		switch env.pkt.Type {
		case PingReq:
			wg.Add(1)
			go func() { conn_p.PingRes(env); wg.Done() }()
		case Store:
			wg.Add(1)
			go func() { conn_p.StoreHandler(env); wg.Done() }()
		default:
			// Acquire write lock and write to cache
			conn_p.mut.Lock()
			conn_p.pcache[env.pkt.RandNum] = env
			conn_p.mut.Unlock()
		}
	}
	wg.Wait()
}

// ReadLoop is supposed to be run as a go routine which can read all the messages comming in
// to the node and send those along, if the TTL has not expired, to the Collector.
func (conn_p *Connector) ReadLoop() {
	for {
		var cmdr []byte
		_, addr_p, err := conn_p.conn.ReadFromUDP(cmdr)
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
			var env Envelope = Envelope{pkt, *addr_p}
			conn_p.rch <- env
		}
	}
}

func (conn_p *Connector) PingRes(env Envelope) {
	env.pkt.Type = PingRes
	conn_p.sch <- env

	// Insert to the NbrMap here.
}

func (conn_p *Connector) StoreHandler(env Envelope) {
	// Insert to the ObjMap here
}
