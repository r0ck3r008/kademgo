// NbrMap package is responsible for,
// 1. Inserting to and deleting from neighbour entries from k-buckets
// 2. Calculating distances between itself and the provided neighbours
// K-Buckets are implemented in the module lru.go
package nbrmap

import (
	"fmt"
	"net"

	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
	"github.com/r0ck3r008/kademgo/writeloop"
)

// NbrMap is a structure that serves as the encapsulation over all the K-Buckets
// and provides the functionality to look up or insert a new neighbour.
type NbrMap struct {
	sz  int
	bkt map[int]*NbrNode
}

// Init is the initiator for the NbrMap and initiates the map of k-buckets.
func (nmap_p *NbrMap) Init() {
	nmap_p = &NbrMap{}
	nmap_p.bkt = make(map[int]*NbrNode)
	nmap_p.sz = 0
}

// Insert is used to insert a new neighbour to its correct k-bucket in NeighbourMap.
// This should be invoked as a go routine.
func (nmap_p *NbrMap) Insert(srchash, dsthash *[utils.HASHSZ]byte, obj *net.IP, wrl_p *writeloop.WriteLoop) {
	var indx int = utils.GetDist(srchash, dsthash)
	nnode_p, ok := nmap_p.bkt[indx]
	if !ok {
		nmap_p.bkt[indx] = nbrnodeinit()
		nnode_p = nmap_p.bkt[indx]
	}

	nnode_p.put(srchash, dsthash, obj, wrl_p)

}

// Get is used to see if a neighbour exists in the NeighbourMap, returns error on failure.
func (nmap_p *NbrMap) Get(srchash, dsthash *[utils.HASHSZ]byte) (*pkt.ObjAddr, error) {
	var indx int = utils.GetDist(srchash, dsthash)
	if node_p, ok := nmap_p.bkt[indx]; ok {
		return node_p.get(dsthash)
	}

	return nil, fmt.Errorf("Not Found!")
}
