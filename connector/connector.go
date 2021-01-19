package connector

import (
	"fmt"
	"net"
	"sync"

	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
)

type Connector struct {
	// conn is the actual UDP listening port
	conn *net.UDPConn
	// rwlock is required to sync pcache
	rwlock *sync.RWMutex
	// sch is the send channel
	sch chan pkt.Envelope
	// rch is the read channel within Connector
	rch chan pkt.Envelope
	// nchan is the channel on which messages to node can be sent
	nchan chan<- pkt.Envelope
	// endchan makes the readloop exit during Connector.DeInit
	endchan chan bool
	// wg is required to wait for rdl and wrl routines to finish
	wg *sync.WaitGroup
	// pcache is the map that stores retured packets for processing
	pcache map[int64]pkt.Envelope
}

func (connP *Connector) Init(nchan chan<- pkt.Envelope) error {
	connP.pcache = make(map[int64]pkt.Envelope)
	connP.rwlock = &sync.RWMutex{}
	connP.wg = &sync.WaitGroup{}
	connP.sch = make(chan pkt.Envelope, 100)
	connP.rch = make(chan pkt.Envelope, 100)
	connP.endchan = make(chan bool)
	connP.nchan = nchan

	var err error
	srvaddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", utils.PORTNUM))
	if err != nil {
		return fmt.Errorf("UDP Resolve: %s", err)
	}
	connP.conn, err = net.ListenUDP("udp", srvaddr)
	if err != nil {
		return fmt.Errorf("UDP Create: %s", err)
	}

	connP.wg.Add(3)
	go func() { connP.readloop(); connP.wg.Done() }()
	go func() { connP.writeloop(); connP.wg.Done() }()
	go func() { connP.collector(); connP.wg.Done() }()

	return nil
}

func (connP *Connector) DeInit() {
	close(connP.sch)
	connP.endchan <- true
	close(connP.rch)
	connP.wg.Wait()
	connP.conn.Close()
}
