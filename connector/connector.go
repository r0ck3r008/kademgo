// connector package is supposed to act like a modified socket connection
// for receiving and sending messages. It can be passed around to functions
// that need to leverage the API.
package connector

import (
	"fmt"
	"net"
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

// DeInit close the UDP bound socket as well as both the send and recv channels.
func (conn_p *Connector) DeInit() {
	// Closing rch and sch would inevitably make ReadLoop/WriteLoop and Collector exit.
	close(conn_p.sch)
	close(conn_p.rch)
	conn_p.conn.Close()
}
