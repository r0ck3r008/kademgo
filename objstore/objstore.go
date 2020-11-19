// Objstore package is responsible for maintaing a prefix tree for the objects
// that the current node possesses. Its provides insertion and checking for existence
package objstore

// ObjNode is the node that represents the Key, Value pair.
// If object is empty slice, it means that it is not a leaf node.
type ObjNode struct {
	hash []byte
	cmap map[byte]*ObjNode
	obj  []byte
}

// ObjStore is an abstraction of the root node of the prefix tree.
type ObjStore struct {
	root *ObjNode
}

// getlvl finds the max character length upto where two hashes are the same.
func getlvl(hash1 []byte, hash2 []byte) int {
	var lvl int = 0
	var mlen int
	if len(hash1) < len(hash2) {
		mlen = len(hash1)
	} else {
		mlen = len(hash2)
	}
	for lvl < mlen {
		if hash1[lvl] != hash2[lvl] {
			break
		}
		lvl++
	}

	return lvl
}

// hash_update function updates the hash map of a particular ObjNode.
func hash_update(obn_p *ObjNode, slice []byte, indx byte, obj []byte, lvl int) {
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
func (obn_p *ObjNode) insert(hash []byte, obj []byte) {
	var lvl int = getlvl(obn_p.hash, hash)
	if lvl > 0 {
		var slice1 []byte = hash[lvl:]
		var slice2 []byte = obn_p.hash[lvl:]
		hash_update(obn_p, slice1, hash[lvl], obj, lvl)
		hash_update(obn_p, slice2, obn_p.hash[lvl], obn_p.obj, lvl)
		obn_p.hash = obn_p.hash[:lvl]
		obn_p.obj = []byte{}
	} else {
		hash_update(obn_p, hash, hash[lvl], obj, 0)
	}
}

// find recursively checks if the given hash is present in the prefix tree.
func (obn_p *ObjNode) find(hash []byte) bool {
	var lvl int = getlvl(obn_p.hash, hash)
	if lvl == len(hash) && len(obn_p.obj) != 0 {
		if lvl == len(obn_p.hash) {
			return true
		} else {
			return false
		}
	} else if node_p, ok := obn_p.cmap[hash[lvl]]; ok {
		return node_p.find(hash[lvl:])
	}

	return false
}

// Init initiates the root in ObjStore.
func (ost_p *ObjStore) Init() {
	var obj []byte
	ost_p.root = &ObjNode{
		[]byte{0},
		make(map[byte]*ObjNode),
		obj,
	}
}

// Insert is an encapsulation over ObjNode.insert.
func (ost_p *ObjStore) Insert(hash []byte, obj []byte) {
	ost_p.root.insert(hash, obj)
}

// Find is an encapsulation over ObjNode.find.
func (ost_p *ObjStore) Find(hash []byte) bool {
	return ost_p.root.find(hash)
}
