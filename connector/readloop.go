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
func (connP *Connector) collector() {
	for env := range connP.rch {
		switch env.Pkt.Type {
		case pkt.PingReq:
			connP.nchan <- env
		case pkt.FindReq:
			connP.nchan <- env
		default:
			// Acquire write lock and write to cache
			connP.rwlock.Lock()
			connP.pcache[env.Pkt.RandNum] = env
			connP.rwlock.Unlock()
		}
	}
}

// ReadLoop is supposed to be run as a go routine which can read all the messages comming in
// to the node and send those along, if the TTL has not expired, to the Collector.
func (connP *Connector) readloop() {
	for {
		select {
		case <-connP.endchan:
			return
		default:
			// Set timeout on read
			if err := connP.conn.SetReadDeadline(time.Now().Add(10 * utils.TIMEOUT)); err != nil {
				fmt.Fprintf(os.Stderr, "Error in setting the read timeout: %s\n", err)
			}

			// Allocate buffer to read into and read
			var cmdr []byte
			_, addrP, err := connP.conn.ReadFromUDP(cmdr)
			if err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					break
				} else {
					fmt.Fprintf(os.Stderr, "Error in reading: %s\n", err)
					os.Exit(1)
				}
			}

			// Unmarshall into packet
			var packet pkt.Packet = pkt.Packet{}
			err = json.Unmarshal(cmdr, &packet)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error in unmarshalling: %s\n", err)
				os.Exit(1)
			}
			var env pkt.Envelope = pkt.Envelope{Pkt: packet, Addr: *addrP}
			// Send out for being handeled
			connP.rch <- env
		}
	}
}
