// kademgo package is responsible for:
// 1. Creating a hash for itself
// 2. Instantiating a nbrmap object
// 3. Instantiating an objstore object
// 4. Providing the API for Kademlia RPCs
package kademgo

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"

	"github.com/r0ck3r008/kademgo/connector"
	"github.com/r0ck3r008/kademgo/nbrmap"
	"github.com/r0ck3r008/kademgo/objmap"
	"github.com/r0ck3r008/kademgo/utils"
)

// KademGo structure encapsulates the UDP listening port, objstore object,
// NbrMap object as well as the hash of the node in question.
type KademGo struct {
	nmap *nbrmap.NbrMap
	ost  *objmap.ObjMap
	hash [utils.HASHSZ]byte
	conn *connector.Connector
	wg   *sync.WaitGroup
}

// Init is the function that initiates the ObjStore, NbrMap, UDP listener
// and as well as forms the random hash for the node.
func (kdm_p *KademGo) Init(addr *string, gway_addr *string) error {
	kdm_p.nmap = &nbrmap.NbrMap{}
	kdm_p.ost = &objmap.ObjMap{}
	kdm_p.conn = &connector.Connector{}

	kdm_p.nmap.Init()
	kdm_p.ost.Init()
	err := kdm_p.conn.Init(addr)
	if err != nil {
		return fmt.Errorf("Error in starting listener: %s\n", err)
	}

	var rnum_str string = strconv.FormatInt(int64(rand.Int()), 10)
	kdm_p.hash = utils.HashStr([]byte(rnum_str))

	kdm_p.wg = &sync.WaitGroup{}
	kdm_p.wg.Add(3)
	go func() { kdm_p.conn.ReadLoop(); kdm_p.wg.Done() }()
	go func() { kdm_p.conn.WriteLoop(); kdm_p.wg.Done() }()
	go func() { kdm_p.conn.Collector(); kdm_p.wg.Done() }()

	kdm_p.conn.FindPeers(&kdm_p.hash, gway_addr)

	return nil
}

// DeInit function waits for all the go routines registered to exit.
func (kdm_p *KademGo) DeInit() {
	kdm_p.conn.DeInit()
	kdm_p.wg.Wait()
}
