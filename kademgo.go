// kademgo implements Kademlia DHT Algorithm for P2P communication in a decentralized
// yet reliable and efficient manner.
package kademgo

import (
	"fmt"
	"os"

	"github.com/r0ck3r008/kademgo/node"
	"github.com/r0ck3r008/kademgo/utils"
)

// KademGo struct is the handle that the user consuming the library would have in order
// to interact with the library.
type KademGo struct {
	node *node.Node
}

// Init initiates the internal structre, node of KademGo.
func (kdm_p *KademGo) Init(addr_p *string, addr_hash *[utils.HASHSZ]byte) {
	var bind_addr string = "127.0.0.1"

	kdm_p.node = &node.Node{}
	if err := kdm_p.node.Init(&bind_addr, addr_p, addr_hash); err != nil {
		fmt.Fprintf(os.Stderr, "Error in initiating node: %s\n", err)
		os.Exit(1)
	}
}

// DeInit calls DeInit on the internal node.
func (kdm_p *KademGo) DeInit() {
	kdm_p.node.DeInit()
}
