// Package nbrmap is responsible for,
// 1. Inserting to and deleting from neighbour entries from k-buckets
// 2. Calculating distances between itself and the provided neighbours
// K-Buckets are implemented in the module lru.go
package nbrmap

import (
	"fmt"
	"net"

	"github.com/r0ck3r008/kademgo/connector"
	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
)

// NbrMap is a structure that serves as the encapsulation over all the K-Buckets
// and provides the functionality to look up or insert a new neighbour.
type NbrMap struct {
	bkt map[int]*NbrNode
}

// Init is the initiator for the NbrMap and initiates the map of k-buckets.
func (nmapP *NbrMap) Init() {
	nmapP.bkt = make(map[int]*NbrNode)
}

// Insert is used to insert a new neighbour to its correct k-bucket in NeighbourMap.
// This should be invoked as a go routine.
func (nmapP *NbrMap) Insert(srchash, dsthash *[utils.HASHSZ]byte, obj *net.IP, connP *connector.Connector) {
	var indx int = utils.GetDist(srchash, dsthash)
	nnodeP, ok := nmapP.bkt[indx]
	if !ok {
		nmapP.bkt[indx] = nbrnodeinit()
		nnodeP = nmapP.bkt[indx]
	}

	nnodeP.put(srchash, dsthash, obj, connP)

}

// Get is used to see if a neighbour exists in the NeighbourMap, returns error on failure.
func (nmapP *NbrMap) Get(srchash, dsthash *[utils.HASHSZ]byte) (*pkt.ObjAddr, error) {
	var indx int = utils.GetDist(srchash, dsthash)
	if nodeP, ok := nmapP.bkt[indx]; ok {
		return nodeP.get(dsthash)
	}

	return nil, fmt.Errorf("not found")
}

// NodeLookup returns the `k' from the bucket closest to the destination, if number of nodes in the
// bucket is lesser than `k', it then searches in buckets other than this until total return number is `k'.
func (nmapP *NbrMap) NodeLookup(srchash, dsthash *[utils.HASHSZ]byte, ret *[]pkt.ObjAddr, sz int) {
	var indx int = utils.GetDist(srchash, dsthash)
	if nodeP, ok := nmapP.bkt[indx]; ok && nodeP.sz == utils.KVAL {
		// A Deep copy of objects
		for i := 0; i < sz; i++ {
			(*ret)[i] = *(nodeP.cvec[i])
		}
	} else {
		// Get neighbours from from earlier buckets since lesser index buckets are closer.
		var i int = 0
		var indxTmp int = indx
		for indxTmp >= 0 && i < sz {
			if nodeP, ok := nmapP.bkt[indxTmp]; ok {
				for j := 0; j < nodeP.sz; j++ {
					(*ret)[i] = *(nodeP.cvec[j])
					i++
				}
			}
			indxTmp--
		}
		// If still not done, get nbrs from rest of the buckets
		indxTmp = indx + 1
		for indxTmp < utils.HASHSZ && i < sz {
			if nodeP, ok := nmapP.bkt[indxTmp]; ok {
				for j := 0; j < nodeP.sz; j++ {
					(*ret)[i] = *(nodeP.cvec[j])
					i++
				}
			}
			indxTmp++
		}

	}
}
