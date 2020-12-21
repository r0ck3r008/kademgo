// node package is responsible for:
// 1. Creating a hash for itself
// 2. Creating a UDPConn on the given bind address.
// 3. Instantiating Reader and WriterLoop objects.
// 4. Initiating Reader, Writer and Collector loops.
package node

import (
	"math/rand"
	"strconv"

	"github.com/r0ck3r008/kademgo/connector"
	"github.com/r0ck3r008/kademgo/nbrmap"
	"github.com/r0ck3r008/kademgo/objmap"
	"github.com/r0ck3r008/kademgo/utils"
)

// Node structure encapsulates the UDP listening port, objstore object,
// NbrMap object as well as the hash of the node in question.
type Node struct {
	hash   [utils.HASHSZ]byte
	conn   *connector.Connector
	omap   *objmap.ObjMap
	nmap   *nbrmap.NbrMap
}

// Init is the function that initiates the ReaderLoop, WriterLoop, UDP listener
// and as well as forms the random hash for the node.
func (node_p *Node) Init(addr *string, gway_addr *string) error {
	var rnum_str string = strconv.FormatInt(int64(rand.Int()), 10)
	node_p.hash = utils.HashStr([]byte(rnum_str))

	node_p.omap = &objmap.ObjMap{}
	node_p.nbrmap = &nbrmap.NbrMap{}

	node_p.omap.Init()
	node_p.nbrmap.Init()

	node_p.conn = &connector.Connector{}
	if err := node_p.conn.Init(addr); err != nil {
		return err
	}

	return nil
}

// DeInit function waits for all the go routines registered to exit.
func (node_p *Node) DeInit() {
	node_p.conn.DeInit()
}
