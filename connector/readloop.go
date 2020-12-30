package connector

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
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
	var flag bool = true
	for flag {
		select {
		case <-conn_p.endchan:
			flag = false
		default:
			// Set timeout on read
			if err := conn_p.conn.SetReadDeadline(time.Now().Add(10 * utils.TIMEOUT)); err != nil {
				fmt.Fprintf(os.Stderr, "Error in setting the read timeout: %s\n", err)
			}

			// Allocate buffer to read into and read
			var cmdr []byte
			_, addr_p, err := conn_p.conn.ReadFromUDP(cmdr)
			if err != nil {
				if !errors.Is(err, os.ErrDeadlineExceeded) {
					fmt.Fprintf(os.Stderr, "Error in reading: %s\n", err)
				}
				break
			}

			// Unmarshall into packet
			var packet pkt.Packet = pkt.Packet{}
			err = json.Unmarshal(cmdr, &packet)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error in unmarshalling: %s\n", err)
				os.Exit(1)
			}
			var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: *addr_p}
			// Send out for being handeled
			conn_p.rch <- env
		}
	}
}
