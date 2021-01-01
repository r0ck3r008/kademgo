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

func (conn_p *Connector) Init(nchan chan<- pkt.Envelope) error {
	conn_p.pcache = make(map[int64]pkt.Envelope)
	conn_p.rwlock = &sync.RWMutex{}
	conn_p.wg = &sync.WaitGroup{}
	conn_p.sch = make(chan pkt.Envelope, 100)
	conn_p.rch = make(chan pkt.Envelope, 100)
	conn_p.endchan = make(chan bool)
	conn_p.nchan = nchan

	var err error
	srvaddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", utils.PORTNUM))
	if err != nil {
		return fmt.Errorf("UDP Resolve: %s", err)
	}
	conn_p.conn, err = net.ListenUDP("udp", srvaddr)
	if err != nil {
		return fmt.Errorf("UDP Create: %s", err)
	}

	conn_p.wg.Add(3)
	go func() { conn_p.readloop(); conn_p.wg.Done() }()
	go func() { conn_p.writeloop(); conn_p.wg.Done() }()
	go func() { conn_p.collector(); conn_p.wg.Done() }()

	return nil
}

func (conn_p *Connector) DeInit() {
	close(conn_p.sch)
	conn_p.endchan <- true
	close(conn_p.rch)
	conn_p.wg.Wait()
	conn_p.conn.Close()
}
