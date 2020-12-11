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
	bkt map[int]*NbrNode
}

// Init is the initiator for the NbrMap and initiates the map of k-buckets.
func (nmap_p *NbrMap) Init() {
	nmap_p = &NbrMap{}
	nmap_p.bkt = make(map[int]*NbrNode)
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

// NodeLookup returns the `k' from the bucket closest to the destination, if number of nodes in the
// bucket is lesser than `k', it then searches in buckets other than this until total return number is `k'.
func (nmap_p *NbrMap) NodeLookup(srchash, dsthash *[utils.HASHSZ]byte, ret []pkt.ObjAddr, sz int) {
	var indx int = utils.GetDist(srchash, dsthash)
	if node_p, ok := nmap_p.bkt[indx]; ok && node_p.sz == utils.KVAL {
		// A Deep copy of objects
		for i := 0; i < sz; i++ {
			ret[i] = *(node_p.cvec[i])
		}
	} else {
		// Get neighbours from from earlier buckets since lesser index buckets are closer.
		var i int = 0
		var indx_tmp int = indx
		for indx_tmp < utils.HASHSZ && i < sz {
			if node_p, ok := nmap_p.bkt[indx_tmp]; ok {
				for j := 0; j < node_p.sz; j++ {
					ret[i] = *(node_p.cvec[j])
					i--
				}
			}
			indx_tmp--
			if indx_tmp == 0 {
				indx_tmp = indx + 1
			}
		}
	}
}
