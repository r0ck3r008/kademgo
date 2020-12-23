// node package is responsible for:
// 1. Creating a hash for itself
// 2. Creating a UDPConn on the given bind address.
// 3. Instantiating Reader and WriterLoop objects.
// 4. Initiating Reader, Writer and Collector loops.
package node

import (
	"math/rand"
	"strconv"
	"sync"

	"github.com/r0ck3r008/kademgo/connector"
	"github.com/r0ck3r008/kademgo/nbrmap"
	"github.com/r0ck3r008/kademgo/objmap"
	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
)

// Node structure encapsulates the UDP listening port, objstore object,
// NbrMap object as well as the hash of the node in question.
type Node struct {
	// hash is the SHA512 hash ID of the node.
	hash  [utils.HASHSZ]byte
	nchan chan pkt.Envelope
	conn  *connector.Connector
	omap  *objmap.ObjMap
	nmap  *nbrmap.NbrMap
	wg    *sync.WaitGroup
}

// Init is the function that initiates the ReaderLoop, WriterLoop, UDP listener
// and as well as forms the random hash for the node.
func (node_p *Node) Init(addr *string, gway_addr *string) error {
	var rnum_str string = strconv.FormatInt(int64(rand.Int()), 10)
	node_p.hash = utils.HashStr([]byte(rnum_str))

	node_p.omap = &objmap.ObjMap{}
	node_p.nmap = &nbrmap.NbrMap{}

	node_p.omap.Init()
	node_p.nmap.Init()

	node_p.conn = &connector.Connector{}
	if err := node_p.conn.Init(addr, node_p.nchan); err != nil {
		return err
	}

	node_p.nchan = make(chan pkt.Envelope)
	node_p.wg = &sync.WaitGroup{}

	node_p.wg.Add(1)
	go func() { node_p.collector(); node_p.wg.Done() }()

	return nil
}

// collector in Node collects all the packets from the readloop which require the
// involvement of NbrMap or ObjMap and appropriately processes them.
func (node_p *Node) collector() {
	wg := sync.WaitGroup{}
	for env := range node_p.nchan {
		switch env.Pkt.Type {
		case pkt.PingRes:
			wg.Add(1)
			go func() { node_p.nmap.Insert(node_p.hash, env.Pkt.Hash, env.Addr.IP, node_p.conn); wg.Done() }()
		case pkt.Store:
			wg.Add(1)
			go func() { node_p.omap.Insert(node_p.hash, env.Pkt.Obj); wg.Done() }()
		}
	}
	wg.Wait()
}

// DeInit function waits for all the go routines registered to exit.
func (node_p *Node) DeInit() {
	node_p.conn.DeInit()
}
