// Node package is responsible for:
// 1. Creating a hash for itself
// 2. Instantiating a nbrmap object
// 3. Instantiating an objstore object
// 4. Providing the API for Kademlia RPCs
// Kademlia API for RPCs is implemented in api.go
package kademgo

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
type KademGo struct {
	nmap *nbrmap.NbrMap
	ost  *objstore.ObjStore
	hash [utils.HASHSZ]byte
	conn *connector.Connector
	wg   *sync.WaitGroup
}

// NodeInit is the function that initiates the ObjStore, NbrMap, UDP listener
// and as well as forms the random hash for the node.
func (kdm_p *KademGo) Init(addr *string) error {
	kdm_p.nmap = nbrmap.NbrMapInit()
	kdm_p.ost = objstore.ObjStoreInit()

	var rnum_str string = strconv.FormatInt(int64(rand.Int()), 10)
	kdm_p.hash = utils.HashStr([]byte(rnum_str))

	conn, err := connector.ConnectorInit(addr)
	if err != nil {
		return fmt.Errorf("Error in starting listener: %s\n", err)
	}
	kdm_p.conn = conn
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
	kdm_p.wg = wg

	return nil
}

func (kdm_p *KademGo) DeInit() {
	kdm_p.wg.Wait()
}
