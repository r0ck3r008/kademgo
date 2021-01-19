// Package objmap stores objects as a mapping of their distance from
// the node. Each object is stored in a vector contained by ObjNode at
// its correct distance bucket. The `key' in the <Key, Value> pair is
// the hash of the object being referred to while the `Value' is the
// address of the peer where it can be found.
package objmap

import (
	"fmt"

	"github.com/r0ck3r008/kademgo/pkt"
	"github.com/r0ck3r008/kademgo/utils"
)

// ObjNode is a single bucket that stores the objects in a vector.
// Each object has its index mapped using its hash within the same node.
type ObjNode struct {
	Nmap map[[utils.HASHSZ]byte]int
	Nvec []*pkt.ObjAddr
}

// ObjMap is the high level mapping of distances from the node to seprate buckets.
type ObjMap struct {
	omap map[int]*ObjNode
}

// Init initialized the ObjMap.
func (omapP *ObjMap) Init() {
	omapP.omap = make(map[int]*ObjNode)
}

// Insert inserts the object accoring to its distance from the node.
func (omapP *ObjMap) Insert(srchash [utils.HASHSZ]byte, obj pkt.ObjAddr) {
	var indx int = utils.GetDist(&srchash, &obj.Hash)
	if nodeP, ok := omapP.omap[indx]; ok {
		if _, ok := nodeP.Nmap[obj.Hash]; !ok {
			nodeP.Nmap[obj.Hash] = len(nodeP.Nvec)
			nodeP.Nvec = append(nodeP.Nvec, &pkt.ObjAddr{Hash: obj.Hash, Addr: obj.Addr})
		}
	} else {
		var nodeP *ObjNode = &ObjNode{Nmap: make(map[[utils.HASHSZ]byte]int),
			Nvec: []*pkt.ObjAddr{&pkt.ObjAddr{Hash: obj.Hash, Addr: obj.Addr}}}
		nodeP.Nmap[obj.Hash] = 0
		omapP.omap[indx] = nodeP
	}
}

// Get fetches the object if it exists.
func (omapP *ObjMap) Get(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte) (*pkt.ObjAddr, error) {
	var indx int = utils.GetDist(srchash, dsthash)
	if nodeP, ok := omapP.omap[indx]; ok {
		if k, ok := nodeP.Nmap[*dsthash]; ok {
			return nodeP.Nvec[k], nil
		}
	}

	return nil, fmt.Errorf("not found")
}

// GetAll returns all the objects that are closer to the given hash from the src hash as a list of
// pointers to a slice of objects.
func (omapP *ObjMap) GetAll(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte) []*[]*pkt.ObjAddr {
	var indx int = utils.GetDist(srchash, dsthash)
	var ret []*[]*pkt.ObjAddr
	for i := 1; i <= indx; i++ {
		if nodeP, ok := omapP.omap[i]; ok {
			ret = append(ret, &nodeP.Nvec)
		}
	}

	return ret
}
