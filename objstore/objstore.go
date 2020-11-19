package objmap

import (
	"fmt"

	"github.com/r0ck3r008/kademgo/utils"
)

type ObjNode struct {
	nmap map[[utils.HASHSZ]byte]int
	nvec []interface{}
}

type ObjMap struct {
	omap map[int]*ObjNode
}

func (omap_p *ObjMap) Init() {
	omap_p.omap = make(map[int]*ObjNode)
}

func (omap_p *ObjMap) Insert(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte, obj interface{}) {
	var indx int = utils.GetDist(srchash, dsthash)
	if node_p, ok := omap_p.omap[indx]; ok {
		if _, ok := node_p.nmap[*dsthash]; !ok {
			omap_p.omap[indx].nmap[*dsthash] = len(omap_p.omap[indx].nvec)
			omap_p.omap[indx].nvec = append(omap_p.omap[indx].nvec, obj)
		}
	} else {
		var node_p *ObjNode = &ObjNode{make(map[[utils.HASHSZ]byte]int), []interface{}{obj}}
		node_p.nmap[*dsthash] = 0
		omap_p.omap[indx] = node_p
	}
}

func (omap_p *ObjMap) Get(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte) (*interface{}, error) {
	var indx int = utils.GetDist(srchash, dsthash)
	if node_p, ok := omap_p.omap[indx]; ok {
		if k, ok := node_p.nmap[*dsthash]; ok {
			return &omap_p.omap[indx].nvec[k], nil
		}
	}

	return nil, fmt.Errorf("Not Found!")
}

func (omap_p *ObjMap) GetAll(srchash *[utils.HASHSZ]byte, dsthash *[utils.HASHSZ]byte) []*[]interface{} {
	var indx int = utils.GetDist(srchash, dsthash)
	var ret []*[]interface{}
	for i := 1; i <= indx; i++ {
		if node_p, ok := omap_p.omap[i]; ok {
			ret = append(ret, &node_p.nvec)
		}
	}

	return ret
}
