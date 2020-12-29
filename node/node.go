// node package is responsible for:
// 1. Creating a hash for itself
// 2. Creating a UDPConn on the given bind address.
// 3. Instantiating Reader and WriterLoop objects.
// 4. Initiating Reader, Writer and Collector loops.
package node

import (
	"math/rand"
	"net"
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
func (node_p *Node) Init(addr *string, gway_addr *string, gway_hash *[utils.HASHSZ]byte) error {
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

	if gway_addr != nil && gway_hash != nil {
		// Insert the gateway node to NbrTable
		var gway_ip net.IP = net.IP([]byte(*gway_addr))
		node_p.nmap.Insert(&node_p.hash, gway_hash, &gway_ip, node_p.conn)
		// Run a lookup on self
		node_p.FindNode(node_p.hash)
	}

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
			go func() { node_p.PingReqHandler(env); wg.Done() }()
		case pkt.FindReq:
			wg.Add(1)
			go func() { node_p.FindReqHandler(env); wg.Done() }()
		}
	}
	wg.Wait()
}

// FindNode is responsible for beginning the process of lookup by
// calling the Connector's FindNode.
func (node_p *Node) FindNode(target [utils.HASHSZ]byte) {
	var ret []pkt.ObjAddr = make([]pkt.ObjAddr, utils.ALPHAVAL)
	// Get First ALPHANUM Nbrs
	node_p.nmap.NodeLookup(&node_p.hash, &target, &ret, utils.ALPHAVAL)

	// Begin the Recursion.
	_ = node_p.conn.FindNodeReq(&node_p.hash, &target, &ret)
}

// DeInit function waits for all the go routines registered to exit.
func (node_p *Node) DeInit() {
	// first closing the nchan on Node.collector helps make sure it exists before
	// we call wait on it.
	close(node_p.nchan)
	node_p.wg.Wait()
	node_p.conn.DeInit()
}
