package connector

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/r0ck3r008/kademgo/pkt"
)

// Collector is intended to be a goroutine that process the received packets in form of Envelope
// struct and caches it in the connector cache based on the identifier.
func (conn_p *Connector) collector() {
	for env := range conn_p.rch {
		switch env.Pkt.Type {
		case pkt.PingReq:
			conn_p.nchan <- env
		case pkt.FindReq:
			conn_p.nchan <- env
		default:
			// Acquire write lock and write to cache
			conn_p.rwlock.Lock()
			conn_p.pcache[env.Pkt.RandNum] = env
			conn_p.rwlock.Unlock()
		}
	}
}

// ReadLoop is supposed to be run as a go routine which can read all the messages comming in
// to the node and send those along, if the TTL has not expired, to the Collector.
func (conn_p *Connector) readloop() {
	for {
		var cmdr []byte
		_, addr_p, err := conn_p.conn.ReadFromUDP(cmdr)
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
		// BUG: This might make the application panic if DeInit on Connector is called while
		// receive channel is being written to with a new packet.
		conn_p.rch <- env
	}
}
