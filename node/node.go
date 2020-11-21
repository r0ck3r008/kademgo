// kademgo package is responsible for:
// 1. Creating a hash for itself
// 2. Instantiating a nbrmap object
// 3. Instantiating an objstore object
// 4. Providing the API for Kademlia RPCs
package node

import (
	"crypto/rand"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/r0ck3r008/kademgo/readloop"
	"github.com/r0ck3r008/kademgo/utils"
	"github.com/r0ck3r008/kademgo/writeloop"
)

// Node structure encapsulates the UDP listening port, objstore object,
// NbrMap object as well as the hash of the node in question.
type Node struct {
	hash   [utils.HASHSZ]byte
	sch    chan utils.Envelope
	pcache map[int64]utils.Envelope
	conn   *net.UDPConn
	wg     *sync.WaitGroup
	rdl    *readloop.ReadLoop
	wrl    *writeloop.WriteLoop
	mut    *sync.Mutex
}

// Init is the function that initiates the ObjStore, NbrMap, UDP listener
// and as well as forms the random hash for the node.
func (node_p *Node) Init(addr *string, gway_addr *string) error {
	node_p.pcache = make(map[int64]utils.Envelope)
	node_p.mut = &sync.Mutex{}
	node_p.wg = &sync.WaitGroup{}
	var rnum_str string = strconv.FormatInt(int64(rand.Int()), 10)
	node_p.hash = utils.HashStr([]byte(rnum_str))
	node_p.sch = make(chan utils.Envelope, 100)

	node_p.rdl = &readloop.ReadLoop{}
	node_p.wrl = &writeloop.WriteLoop{}
	node_p.rdl.Init(rdl_p.mut, &node_p.pcache, node_p.sch)
	node_p.wrl.Init(rdl_p.mut, &node_p.pcache, node_p.sch)

	var err error
	node_p.conn, err = net.ListenUDP("conn", &net.UDPAddr{IP: []byte(*addr), Port: utils.PORTNUM, Zone: ""})
	if err != nil {
		return fmt.Errorf("UDP Create: %s", err)
	}
	node_p.wg.Add(3)
	go func() { node_p.rdl.ReadLoop(node_p.conn); kdm_p.wg.Done() }()
	go func() { node_p.wrl.WriteLoop(node_p.conn); kdm_p.wg.Done() }()
	go func() { node_p.rdl.Collector(); kdm_p.wg.Done() }()

	return nil
}

// DeInit function waits for all the go routines registered to exit.
func (node_p *Node) DeInit() {
	close(node_p.sch)
	node_p.rdl.DeInit()
	node_p.wrl.DeInit()
	node_p.wg.Wait()
}
