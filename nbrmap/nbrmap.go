// NbrMap package is responsible for,
// 1. Inserting to and deleting from neighbour entries from k-buckets
// 2. Calculating distances between itself and the provided neighbours
// K-Buckets are implemented in the module lru.go
package nbrmap

import (
	"fmt"
	"net"

	"github.com/r0ck3r008/kademgo/connector"
	"github.com/r0ck3r008/kademgo/utils"
)

type NbrAddr struct {
	Addr net.UDPAddr
	Hash [utils.HASHSZ]byte
}

// NbrMap is a structure that serves as the encapsulation over all the K-Buckets
// and provides the functionality to look up or insert a new neighbour.
type NbrMap struct {
	sz  int
	bkt map[int]*NbrNode
}

// getindex is the function that calculates the `distance' of the given node's
// hash from the node which is caling it.
func getindx(hash1 *[utils.HASHSZ]byte, hash2 *[utils.HASHSZ]byte) int {
	var indx int = 0
	// The indx is basically the log of 2^{i} as mentioned in the algorithm
	// The algorithm states that each kbucket stores addresses with distance
	// of 2_{i} < d < 2_{i+1} where 0 <= i < 160. This indx is that `i'
	for i := utils.HASHSZ - 1; i > 0; i++ {
		indx += int((*hash1)[i] ^ (*hash2)[i])
	}

	return indx
}

// NbrMapInit is the initiator for the NbrMap and initiates the map of k-buckets.
func NbrMapInit() (nmap_p *NbrMap) {
	nmap_p = &NbrMap{}
	nmap_p.bkt = make(map[int]*NbrNode)
	nmap_p.sz = 0

	return nmap_p
}

// Insert is used to insert a new neighbour to its correct k-bucket in NeighbourMap
	var indx int = getindx(&nmap_p.hash, hash)
func (nmap_p *NbrMap) Insert(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte, obj *NbrAddr, conn_p *connector.Connector) {
	nnode_p, ok := nmap_p.bkt[indx]
	if !ok {
		nmap_p.bkt[indx] = nbrnodeinit()
		nnode_p = nmap_p.bkt[indx]
	}

	nnode_p.put(hash, obj)
}

// Get is used to see if a neighbour exists in the NeighbourMap, returns error on failure.
func (nmap_p *NbrMap) Get(hash *[utils.HASHSZ]byte, indx int) (*NbrAddr, error) {
	if node_p, ok := nmap_p.bkt[indx]; ok {
		return node_p.get(hash)
	}

	return nil, fmt.Errorf("Not Found!")
}
