// connector package is supposed to act like a modified socket connection
// for receiving and sending messages. It can be passed around to functions
// that need to leverage the API.
package connector

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	sync "sync"

	"github.com/r0ck3r008/kademgo/utils"
)

// Envelope is an encapsulation which would be passed around in go channels.
type Envelope struct {
	pkt  Pkt
	addr net.UDPAddr
}

// Connector type that stores all the channel, wait mutex and packet cache
// and is an required element before any function can use the API.
type Connector struct {
	conn   *net.UDPConn
	pcache map[int64]Envelope
	mut    *sync.Mutex
	sch    chan Envelope
	rch    chan Envelope
}

// Init sets up the UDP listening socket, send and recv channels, mutex and the packet cache map.
func (conn_p *Connector) Init(addr *string) error {
	var err error
	conn_p.conn, err = net.ListenUDP("conn", &net.UDPAddr{IP: []byte(*addr), Port: utils.PORTNUM, Zone: ""})
	if err != nil {
		return fmt.Errorf("UDP Create: %s", err)
	}
	conn_p.sch = make(chan Envelope, 100)
	conn_p.rch = make(chan Envelope, 100)
	conn_p.mut = &sync.Mutex{}

	return nil
}

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
			close(conn_p.rch)
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
			conn_p.sch <- env
		}
	}
}

// WriteLoop is supposed to be run as a goroutine which takes all the packets that need to be sent
// from the node and send them asynchronously to the desired destinations.
func (conn_p *Connector) WriteLoop() {
	for env := range conn_p.sch {
		cmds, err := json.Marshal(env.pkt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Marshalling: %s\n", err)
			break
		}
		if _, err := conn_p.conn.WriteToUDP(cmds, &env.addr); err != nil {
			fmt.Fprintf(os.Stderr, "Error in writing: %s\n", err)
			break
		}
	}
}
