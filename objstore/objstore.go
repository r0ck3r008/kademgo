// Objstore package is responsible for maintaing a prefix tree for the objects
// that the current node possesses. Its provides insertion and checking for existence
package objstore

// ObjNode embodies each node in the prefix tree for the ObjStore.
type ObjNode struct {
	hash []byte
	cmap map[byte]*ObjNode
	obj  interface{}
}

// ObjStore is the encapsulation over the root node of the prefix tree.
type ObjStore struct {
	root *ObjNode
}

// getlvl finds the max character length upto where two hashes are the same.
func getlvl(hash1 []byte, hash2 []byte) int {
	var lvl int = -1
	var mlen int
	if len(hash1) < len(hash2) {
		mlen = len(hash1)
	} else {
		mlen = len(hash2)
	}
	for lvl < mlen {
		lvl++
		if hash1[lvl] != hash2[lvl] {
			break
		}
	}

	return lvl
}

// hash_update function updates the hash map of a particular ObjNode.
func hash_update(obn_p *ObjNode, slice []byte, indx byte, obj interface{}, lvl int) {
	if node_p, ok := obn_p.cmap[indx]; ok {
		node_p.insert(slice, obj)
	} else {
		obn_p.cmap[indx] = &ObjNode{
			slice,
			make(map[byte]*ObjNode),
			obj,
		}
	}
}

// insert recursively inserts an object according to its hash into the prefix tree.
func (obn_p *ObjNode) insert(hash []byte, obj interface{}) {
	var lvl int = getlvl(obn_p.hash, hash)
	if lvl >= 0 {
		var slice1 []byte = hash[lvl:]
		var slice2 []byte = obn_p.hash[lvl:]
		var indx byte = hash[lvl]
		hash_update(obn_p, slice1, indx, obj, lvl)
		hash_update(obn_p, slice2, indx, obn_p.obj, lvl)
		obn_p.hash = obn_p.hash[:lvl]
		obn_p.obj = nil
	}
}

// find recursively checks if the given hash is present in the prefix tree.
func (obn_p *ObjNode) find(hash []byte) bool {
	var lvl int = getlvl(obn_p.hash, hash)
	if lvl == utils.KVAL-1 {
		return true
	} else if lvl >= 0 {
		var indx byte = hash[lvl]
		if node_p, ok := obn_p.cmap[indx]; ok {
			var slice []byte = hash[lvl:]
			return node_p.find(slice)
		}
	}

	return false
}

// ObjStoreInit initiates the root in ObjStore.
func ObjStoreInit() (ost_p *ObjStore) {
	ost_p = &ObjStore{}
	ost_p.root = &ObjNode{
		[]byte{0},
		make(map[byte]*ObjNode),
		nil,
	}

	return ost_p
}

// Insert is an encapsulation over ObjNode.insert.
func (ost_p *ObjStore) Insert(hash []byte, obj interface{}) {
	ost_p.root.insert(hash, obj)
}

// Find is an encapsulation over ObjNode.find.
func (ost_p *ObjStore) Find(hash []byte) bool {
	return ost_p.root.find(hash)
}
