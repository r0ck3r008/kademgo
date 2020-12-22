package connector

import (
	"fmt"
	"sync"
	"net"

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
	// wg is required to wait for rdl and wrl routines to finish
	wg *sync.WaitGroup
	// pcache is the map that stores retured packets for processing
	pcache map[int64]pkt.Envelope
}

func (conn_p *Connector) Init(addr *string) error {
	conn_p.pcache = make(map[int64]pkt.Envelope)
	conn_p.rwlock = &sync.RWMutex{}
	conn_p.wg = &sync.WaitGroup{}
	conn_p.sch = make(chan pkt.Envelope, 100)
	conn_p.rch = make(chan pkt.Envelope, 100)

	var err error
	conn_p.conn, err = net.ListenUDP("conn", &net.UDPAddr{IP: []byte(*addr), Port: utils.PORTNUM, Zone: ""})
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
	close(conn_p.rch)
	conn_p.wg.Wait()
	conn_p.conn.Close()
}
