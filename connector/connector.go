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
	// rdl is the readloop object
	rdl *readloop
	// wrl is the writeloop object
	wrl *writeloop
	// rwlock is required to sync pcache
	rwlock *sync.RWMutex
	// sch is the send channel
	sch chan pkt.Envelope
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

	var err error
	conn_p.conn, err = net.ListenUDP("conn", &net.UDPAddr{IP: []byte(*addr), Port: utils.PORTNUM, Zone: ""})
	if err != nil {
		return fmt.Errorf("UDP Create: %s", err)
	}

	conn_p.rdl = &readloop{}
	conn_p.wrl = &writeloop{}
	conn_p.rdl.init(conn_p.rwlock, &conn_p.pcache, conn_p.sch)
	conn_p.wrl.init(conn_p.rwlock, &conn_p.pcache, conn_p.sch)

	conn_p.wg.Add(3)
	go func() { conn_p.rdl.readloop(conn_p.conn); conn_p.wg.Done() }()
	go func() { conn_p.wrl.writeloop(conn_p.conn); conn_p.wg.Done() }()
	go func() { conn_p.rdl.collector(); conn_p.wg.Done() }()

	return nil
}

func (conn_p *Connector) DeInit() {
	conn_p.wg.Wait()
	conn_p.conn.Close()
}
