package connector

import (
	"encoding/json"
	"fmt"
	"os"
)

// Collector is intended to be a goroutine that process the received packets in form of Envelope
// struct and caches it in the connector cache based on the identifier.
func (conn_p *Connector) Collector() {
	for env := range conn_p.sch {
		// Acquire write lock and write to cache
		conn_p.mut.Lock()
		conn_p.pcache[env.pkt.RandNum] = env
		conn_p.mut.Unlock()
	}
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
}
