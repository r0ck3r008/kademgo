package nbrmap

import (
	"fmt"

	"github.com/r0ck3r008/kademgo/utils"
)

// access maps the address of a neighbor back to its saved hash.
type access struct {
	obj  interface{}
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
func (cache_p *NbrNode) put(hash *[utils.HASHSZ]byte, obj interface{}) {
	if indx, ok := cache_p.cmap[*hash]; ok && (indx != len(cache_p.cvec)-1) {
		// Found! Now remove it from where ever it is
		cache_p.cvec = append(cache_p.cvec[:indx], cache_p.cvec[indx+1:]...)
	} else {
		// Not Found
		// if length of cvec is >= KVAL, remove the first element from front
		// of cvec and cmap
		veclen := len(cache_p.cvec)
		if veclen == utils.KVAL {
			var old_p *access
			old_p, cache_p.cvec = cache_p.cvec[1], cache_p.cvec[1:]
			delete(cache_p.cmap, old_p.hash)
		}
		cache_p.cmap[*hash] = veclen
	}
	// Append to the end now
	cache_p.cvec = append(cache_p.cvec, &access{obj, *hash})
}

// get fetches the neighbour if it exists in the cache, returns error on faliure.
func (cache_p *NbrNode) get(hash *[utils.HASHSZ]byte) (*interface{}, error) {
	if indx, ok := cache_p.cmap[*hash]; ok {
		return &cache_p.cvec[indx].obj, nil
	}

	return nil, fmt.Errorf("Not Found!")
}
