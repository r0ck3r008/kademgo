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
func nbrnodeinit() (cache_p *NbrNode) {
	cache_p = &NbrNode{}
	cache_p.cmap = make(map[[utils.HASHSZ]byte]int)
	cache_p.cvec = make([]*pkt.ObjAddr, utils.KVAL)
	cache_p.sz = 0

	return cache_p
}

// put is used to insert a new neighbour into the cached list of neighbours.
// It uses LRU eviction policy based on an irresponsive last contacted neighbour
// in case of a filled bucket.
func (cache_p *NbrNode) put(srchash, dsthash *[utils.HASHSZ]byte, obj *net.IP, conn_p *connector.Connector) {
	if indx, ok := cache_p.cmap[*dsthash]; ok && (indx != len(cache_p.cvec)-1) {
		// Found! Now remove it from where ever it is and push to the back
		cache_p.cvec = append(cache_p.cvec[:indx], cache_p.cvec[indx+1:]...)
		cache_p.cvec = append(cache_p.cvec, &pkt.ObjAddr{Hash: *dsthash, Addr: *obj})
	} else {
		// Not Found
		var old_p *pkt.ObjAddr = cache_p.cvec[0]

		// if length of cvec is == KVAL, remove the first element from front
		// of cvec and cmap if ping of the least recently used fails
		if cache_p.sz == utils.KVAL {
			cache_p.cvec = cache_p.cvec[1:]
			if conn_p.PingReq(srchash, &old_p.Addr) {
				// If ping succeedes, add the old one to the back
				cache_p.cvec = append(cache_p.cvec, old_p)
				return
			} else {
				// If ping fails delete the lease used one
				delete(cache_p.cmap, old_p.Hash)
				cache_p.sz--
			}
		}
		// If it reaches here, append the new one and increase the Sz by 1
		cache_p.cvec = append(cache_p.cvec, &pkt.ObjAddr{Hash: *dsthash, Addr: *obj})
		cache_p.cmap[*dsthash] = cache_p.sz
		cache_p.sz++
	}
}

// get fetches the neighbour if it exists in the cache, returns error on faliure.
func (cache_p *NbrNode) get(hash *[utils.HASHSZ]byte) (*pkt.ObjAddr, error) {
	if indx, ok := cache_p.cmap[*hash]; ok {
		return cache_p.cvec[indx], nil
	}

	return nil, fmt.Errorf("Not Found!")
}
