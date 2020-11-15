package nbrmap

import (
	"fmt"

	utils "github.com/r0ck3r008/kademgo/utils"
)

type NbrMap struct {
	sz   int
	hash [utils.HASHSZ]byte
	bkt  map[int]*NbrNode
}

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

func (nmap_p *NbrMap) Init() {
	nmap_p.bkt = make(map[int]*NbrNode)
	nmap_p.sz = 0
}

func (nmap_p *NbrMap) Insert(hash *[utils.HASHSZ]byte, obj interface{}) {
	var indx int = getindx(&nmap_p.hash, hash)
	nnode_p, ok := nmap_p.bkt[indx]
	if !ok {
		var nnode_tmp *NbrNode = &NbrNode{}
		nnode_tmp.init()
		nmap_p.bkt[indx] = nnode_tmp
		nnode_p = nnode_tmp
	}

	nnode_p.put(hash, obj)
}

func (nmap_p *NbrMap) Get(hash *[utils.HASHSZ]byte) (*interface{}, error) {
	var indx int = getindx(&nmap_p.hash, hash)
	if node_p, ok := nmap_p.bkt[indx]; ok {
		return node_p.get(hash)
	}

	return nil, fmt.Errorf("Not Found!")
}
