// Node package is responsible for:
// 1. Creating a hash for itself
// 2. Instantiating a nbrmap object
// 3. Instantiating an objstore object
// 4. Providing the API for Kademlia RPCs
// Kademlia API for RPCs is implemented in api.go
package node

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"

	"github.com/r0ck3r008/kademgo/connector"
	"github.com/r0ck3r008/kademgo/nbrmap"
	"github.com/r0ck3r008/kademgo/objstore"
	"github.com/r0ck3r008/kademgo/utils"
)

// Node structure encapsulates the UDP listening port, objstore object,
// NbrMap object as well as the hash of the node in question.
type Node struct {
	nmap *nbrmap.NbrMap
	ost  *objstore.ObjStore
	hash [utils.HASHSZ]byte
	conn *connector.Connector
	wg   *sync.WaitGroup
}

// NodeInit is the function that initiates the ObjStore, NbrMap, UDP listener
// and as well as forms the random hash for the node.
func NodeInit(addr *string) (*Node, error) {
	node_p := &Node{}
	node_p.nmap = nbrmap.NbrMapInit()
	node_p.ost = objstore.ObjStoreInit()

	var rnum_str string = strconv.FormatInt(int64(rand.Int()), 10)
	node_p.hash = utils.HashStr([]byte(rnum_str))

	conn, err := connector.ConnectorInit(addr)
	if err != nil {
		return nil, fmt.Errorf("Error in starting listener: %s\n", err)
	}
	node_p.conn = conn
	var wg = &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		conn.ReadLoop()
		wg.Done()
	}()
	go func() {
		conn.WriteLoop()
		wg.Done()
	}()
	go func() {
		conn.Collector()
		wg.Done()
	}()
	node_p.wg = wg

	return node_p, nil
}

func (node_p *Node) DeInit() {
	node_p.wg.Wait()
}
