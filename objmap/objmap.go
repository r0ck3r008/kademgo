// objmap package stores objects as a mapping of their distance from
// the node. Each object is stored in a vector contained by ObjNode at
// its correct distance bucket. The `key' in the <Key, Value> pair is
// the hash of the object being referred to while the `Value' is the
// address of the peer where it can be found.
package objmap

import (
	"fmt"
	"net"

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

// Init initialized the ObjMap
func (omap_p *ObjMap) Init() {
	omap_p.omap = make(map[int]*ObjNode)
}

// Insert inserts the object accoring to its distance from the node.
func (omap_p *ObjMap) Insert(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte, obj *net.IP) {
	var indx int = utils.GetDist(srchash, dsthash)
	if node_p, ok := omap_p.omap[indx]; ok {
		if _, ok := node_p.Nmap[*dsthash]; !ok {
			omap_p.omap[indx].Nmap[*dsthash] = len(omap_p.omap[indx].Nvec)
			omap_p.omap[indx].Nvec = append(omap_p.omap[indx].Nvec, &pkt.ObjAddr{Hash: *dsthash, Addr: *obj})
		}
	} else {
		var node_p *ObjNode = &ObjNode{Nmap: make(map[[utils.HASHSZ]byte]int),
			Nvec: []*pkt.ObjAddr{&pkt.ObjAddr{Hash: *dsthash, Addr: *obj}}}
		node_p.Nmap[*dsthash] = 0
		omap_p.omap[indx] = node_p
	}
}

// Get fetches the object if it exists.
func (omap_p *ObjMap) Get(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte) (*pkt.ObjAddr, error) {
	var indx int = utils.GetDist(srchash, dsthash)
	if node_p, ok := omap_p.omap[indx]; ok {
		if k, ok := node_p.Nmap[*dsthash]; ok {
			return omap_p.omap[indx].Nvec[k], nil
		}
	}

	return nil, fmt.Errorf("Not Found!")
}

// GetAll returns all the objects that are closer to the given hash from the src hash as a list of
// pointers to a slice of objects.
func (omap_p *ObjMap) GetAll(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte) []*[]*pkt.ObjAddr {
	var indx int = utils.GetDist(srchash, dsthash)
	var ret []*[]*pkt.ObjAddr
	for i := 1; i <= indx; i++ {
		if node_p, ok := omap_p.omap[i]; ok {
			ret = append(ret, &node_p.Nvec)
		}
	}

	return ret
}
