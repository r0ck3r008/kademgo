package nbrmap

import (
	"fmt"
	"net"

	"github.com/r0ck3r008/kademgo/connector"
	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
)

// NbrNode serves as the LRU cache for the Neighbours and represents a particular
// K-Bucket within the larger NeighbourMap
type NbrNode struct {
	cmap map[[utils.HASHSZ]byte]int
	cvec []*pkt.ObjAddr
	sz   int
}

// nbrnodeinit function initiates a NbrNode.
func nbrnodeinit() (cacheP *NbrNode) {
	cacheP = &NbrNode{}
	cacheP.cmap = make(map[[utils.HASHSZ]byte]int)
	cacheP.cvec = make([]*pkt.ObjAddr, utils.KVAL)
	cacheP.sz = 0

	return cacheP
}

// put is used to insert a new neighbour into the cached list of neighbours.
// It uses LRU eviction policy based on an irresponsive last contacted neighbour
// in case of a filled bucket.
func (cacheP *NbrNode) put(srchash, dsthash *[utils.HASHSZ]byte, obj *net.IP, connP *connector.Connector) {
	if indx, ok := cacheP.cmap[*dsthash]; ok && (indx != len(cacheP.cvec)-1) {
		// Found! Now remove it from where ever it is and push to the back
		cacheP.cvec = append(cacheP.cvec[:indx], cacheP.cvec[indx+1:]...)
		cacheP.cvec = append(cacheP.cvec, &pkt.ObjAddr{Hash: *dsthash, Addr: *obj})
	} else {
		// Not Found
		var oldP *pkt.ObjAddr = cacheP.cvec[0]

		// if length of cvec is == KVAL, remove the first element from front
		// of cvec and cmap if ping of the least recently used fails
		if cacheP.sz == utils.KVAL {
			cacheP.cvec = cacheP.cvec[1:]
			if connP.PingReq(srchash, &oldP.Addr) {
				// If ping succeedes, add the old one to the back
				cacheP.cvec = append(cacheP.cvec, oldP)
				return
			} else {
				// If ping fails delete the lease used one
				delete(cacheP.cmap, oldP.Hash)
				cacheP.sz--
			}
		}
		// If it reaches here, append the new one and increase the Sz by 1
		cacheP.cvec = append(cacheP.cvec, &pkt.ObjAddr{Hash: *dsthash, Addr: *obj})
		cacheP.cmap[*dsthash] = cacheP.sz
		cacheP.sz++
	}
}

// get fetches the neighbour if it exists in the cache, returns error on faliure.
func (cacheP *NbrNode) get(hash *[utils.HASHSZ]byte) (*pkt.ObjAddr, error) {
	if indx, ok := cacheP.cmap[*hash]; ok {
		return cacheP.cvec[indx], nil
	}

	return nil, fmt.Errorf("not found")
}
