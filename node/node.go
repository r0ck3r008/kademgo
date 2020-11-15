// This package is responsible for:
// 1. Creating a hash for itself
// 2. Instantiating a nbrmap object
// 3. Instantiating an objstore object
// 4. Providing the API for Kademlia RPCs

package node

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/r0ck3r008/kademgo/nbrmap"
	"github.com/r0ck3r008/kademgo/objstore"
	"github.com/r0ck3r008/kademgo/utils"
)

type Node struct {
	nmap *nbrmap.NbrMap
	ost  *objstore.ObjStore
	hash [utils.HASHSZ]byte
	conn *net.UDPConn
}

func NodeInit(addr *string, port int) (*Node, error) {
	node_p := &Node{}
	node_p.nmap = nbrmap.NbrMapInit()
	node_p.ost = objstore.ObjStoreInit()

	rand.Seed(time.Now().UnixNano())
	var rnum_str string = strconv.FormatInt(int64(rand.Int()), 10)
	node_p.hash = utils.HashStr([]byte(rnum_str))

	conn_p, err := net.ListenUDP("conn", &net.UDPAddr{IP: []byte(*addr), Port: port, Zone: ""})
	if err != nil {
		return nil, fmt.Errorf("UDP Create: %s", err)
	}

	node_p.conn = conn_p

	return node_p, nil
}

func (node_p *Node) SrvLoop() {
	for {
		var cmdr [512]byte
		_, _, err := node_p.conn.ReadFromUDP(cmdr[:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in reading: %s\n", err)
			os.Exit(1)
		}
	}
}

func (node_p *Node) DeInit() {
	node_p.conn.Close()
}
