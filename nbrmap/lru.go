package nbrmap

import (
	"fmt"

	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
	"github.com/r0ck3r008/kademgo/writeloop"
)

// access maps the address of a neighbor back to its saved hash.
type access struct {
	obj  *pkt.NbrAddr
	hash [utils.HASHSZ]byte
}

// NbrNode serves as the LRU cache for the Neighbours and represents a particular
// K-Bucket within the larger NeighbourMap
type NbrNode struct {
	cmap map[[utils.HASHSZ]byte]int
	cvec []*access
	sz   int
}

// nbrnodeinit function initiates a NbrNode.
func nbrnodeinit() (cache_p *NbrNode) {
	cache_p = &NbrNode{}
	cache_p.cmap = make(map[[utils.HASHSZ]byte]int)
	cache_p.cvec = make([]*access, utils.KVAL)
	cache_p.sz = 0

	return cache_p
}

// put is used to insert a new neighbour into the cached list of neighbours.
// It uses LRU eviction policy based on an irresponsive last contacted neighbour
// in case of a filled bucket.
func (cache_p *NbrNode) put(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte, obj *pkt.NbrAddr, wrl_p *writeloop.WriteLoop) {
	if indx, ok := cache_p.cmap[*dsthash]; ok && (indx != len(cache_p.cvec)-1) {
		// Found! Now remove it from where ever it is and push to the back
		cache_p.cvec = append(cache_p.cvec[:indx], cache_p.cvec[indx+1:]...)
		var addr pkt.NbrAddr = *obj
		cache_p.cvec = append(cache_p.cvec, &access{&addr, *dsthash})
	} else {
		// Not Found
		// if length of cvec is == KVAL, remove the first element from front
		// of cvec and cmap if ping of the least recently used fails
		veclen := len(cache_p.cvec)
		if veclen == utils.KVAL {
			// Remove the lease recently used
			var old_p *access
			old_p, cache_p.cvec = cache_p.cvec[0], cache_p.cvec[1:]

			if !wrl_p.Ping(srchash, &old_p.obj.Addr) {
				// If ping fails, add the new one
				delete(cache_p.cmap, old_p.hash)
				var addr pkt.NbrAddr = *obj
				cache_p.cvec = append(cache_p.cvec, &access{&addr, *dsthash})
				cache_p.cmap[*dsthash] = veclen
			} else {
				// If ping succeedes, add the old one to back
				cache_p.cvec = append(cache_p.cvec, old_p)
			}
		}
	}
}

// get fetches the neighbour if it exists in the cache, returns error on faliure.
func (cache_p *NbrNode) get(hash *[utils.HASHSZ]byte) (*pkt.NbrAddr, error) {
	if indx, ok := cache_p.cmap[*hash]; ok {
		return cache_p.cvec[indx].obj, nil
	}

	return nil, fmt.Errorf("Not Found!")
}
