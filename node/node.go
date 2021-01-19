// Package node is responsible for:
// 1. Creating a hash for itself
// 2. Creating a UDPConn on the given bind address.
// 3. Instantiating Reader and WriterLoop objects.
// 4. Initiating Reader, Writer and Collector loops.
package node

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
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
func (nodeP *Node) Init(gwayAddr *string, gwayHash *[utils.HASHSZ]byte) error {
	var rnumStr string = strconv.FormatInt(int64(rand.Int()), 10)
	nodeP.hash = utils.HashStr([]byte(rnumStr))

	nodeP.omap = &objmap.ObjMap{}
	nodeP.nmap = &nbrmap.NbrMap{}

	nodeP.omap.Init()
	nodeP.nmap.Init()

	nodeP.conn = &connector.Connector{}
	if err := nodeP.conn.Init(nodeP.nchan); err != nil {
		return err
	}

	nodeP.nchan = make(chan pkt.Envelope)
	nodeP.wg = &sync.WaitGroup{}

	nodeP.wg.Add(1)
	go func() { nodeP.collector(); nodeP.wg.Done() }()

	if gwayAddr != nil && gwayHash != nil {
		// Insert the gateway node to NbrTable
		var gwayIP net.IP = net.IP([]byte(*gwayAddr))
		nodeP.nmap.Insert(&nodeP.hash, gwayHash, &gwayIP, nodeP.conn)
		// Run a lookup on self
		nodeP.FindNode(nodeP.hash)
	}

	fmt.Fprintf(os.Stdout, "[!] Node successfully initiated with hash: \n%s\n", hex.EncodeToString(nodeP.hash[:]))

	return nil
}

// GetHash copies the hash of the node to provided buffer
func (nodeP *Node) GetHash(buf *[utils.HASHSZ]byte) {
	*buf = nodeP.hash
}

// collector in Node collects all the packets from the readloop which require the
// involvement of NbrMap or ObjMap and appropriately processes them.
func (nodeP *Node) collector() {
	wg := sync.WaitGroup{}
	for env := range nodeP.nchan {
		switch env.Pkt.Type {
		case pkt.PingRes:
			wg.Add(1)
			go func(env pkt.Envelope) { nodeP.PingReqHandler(env); wg.Done() }(env)
		case pkt.FindReq:
			wg.Add(1)
			go func(env pkt.Envelope) { nodeP.FindReqHandler(env); wg.Done() }(env)
		}
	}
	wg.Wait()
}

// FindNode is responsible for beginning the process of lookup by
// calling the Connector's FindNode.
func (nodeP *Node) FindNode(target [utils.HASHSZ]byte) {
	var ret []pkt.ObjAddr = make([]pkt.ObjAddr, utils.ALPHAVAL)
	// Get First ALPHANUM Nbrs
	nodeP.nmap.NodeLookup(&nodeP.hash, &target, &ret, utils.ALPHAVAL)

	// Begin the Recursion.
	_ = nodeP.conn.FindNodeReq(&nodeP.hash, &target, &ret)
}

// DeInit function waits for all the go routines registered to exit.
func (nodeP *Node) DeInit() {
	// first closing the nchan on Node.collector helps make sure it exists before
	// we call wait on it.
	close(nodeP.nchan)
	nodeP.wg.Wait()
	nodeP.conn.DeInit()
}
