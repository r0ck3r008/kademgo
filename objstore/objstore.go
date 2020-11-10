package objstore

import utils "github.com/r0ck3r008/kademgo/utils"

type ObjNode struct {
	hash []byte
	cmap map[byte]*ObjNode
	obj  interface{}
}

type ObjStore struct {
	root *ObjNode
}

func getlvl(hash1 []byte, hash2 []byte) int {
	var lvl int = -1
	for lvl < utils.KVAL {
		lvl++
		if (*hash1)[lvl] != (*hash2)[lvl] {
			break
		}
	}

	return lvl
}

func hash_update(obn_p *ObjNode, slice []byte, indx byte, obj interface{}, lvl int) {
	if node_p, ok := obn_p.cmap[indx]; ok {
		node_p.insert(&slice, obj)
	} else {
		obn_p.cmap[indx] = &ObjNode{
			slice,
			make(map[byte]*ObjNode),
			obj,
		}
	}
}

func (obn_p *ObjNode) insert(hash []byte, obj interface{}) {
	var lvl int = getlvl(obn_p.hash, hash)
	if lvl >= 0 {
		var slice1 string = (*hash)[lvl:]
		var slice2 string = obn_p.hash[lvl:]
		var indx byte = (*hash)[lvl]
		hash_update(obn_p, slice1, indx, obj, lvl)
		hash_update(obn_p, slice2, indx, obn_p.obj, lvl)
		obn_p.hash = obn_p.hash[:lvl]
		obn_p.obj = nil
	}
}

func (obn_p *ObjNode) find(hash []byte) bool {
	var lvl int = getlvl(obn_p.hash, hash)
	if lvl == utils.KVAL-1 {
		return true
	} else if lvl >= 0 {
		var indx byte = (*hash)[lvl]
		var indx byte = hash[lvl]
		if node_p, ok := obn_p.cmap[indx]; ok {
			var slice string = (*hash)[lvl:]
			return node_p.find(&slice)
			var slice []byte = hash[lvl:]
			return node_p.find(slice)
		}
	}

	return false
}

func (ost_p *ObjStore) Init() {
	ost_p.root = &ObjNode{
		[]byte{0},
		make(map[byte]*ObjNode),
		nil,
	}
}

func (ost_p *ObjStore) Insert(hash []byte, obj interface{}) {
	ost_p.root.insert(hash, obj)
}

func (ost_p *ObjStore) Find(hash []byte) bool {
	return ost_p.root.find(hash)
}
