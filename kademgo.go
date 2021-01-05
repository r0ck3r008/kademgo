// kademgo implements Kademlia DHT Algorithm for P2P communication in a decentralized
// yet reliable and efficient manner.
package kademgo

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/r0ck3r008/kademgo/node"
	"github.com/r0ck3r008/kademgo/protos"
	"github.com/r0ck3r008/kademgo/utils"
	"google.golang.org/grpc"
)

// KademGo struct is the handle that the user consuming the library would have in order
// to interact with the library.
type KademGo struct {
	node *node.Node
}

// Init initiates the internal structre, node of KademGo.
func (kdm_p *KademGo) Init(addr_p *string, addr_hash *[utils.HASHSZ]byte) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	kdm_p.node = &node.Node{}
	if err := kdm_p.node.Init(addr_p, addr_hash); err != nil {
		fmt.Fprintf(os.Stderr, "Error in initiating node: %s\n", err)
		os.Exit(1)
	}

	gs := grpc.NewServer()
	protos.RegisterKademgoServer(gs, kdm_p)
	// create a TCP socket for inbound server connections
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", utils.GRPCPORTNUM))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Kademgo: Listener: %s\n", err)
		os.Exit(1)
	}

	// listen for requests
	gs.Serve(l)
}

// GetHash RPC call is to get the hash from the node
func (kdm_p *KademGo) GetHash(ctx context.Context, req *protos.Request) (*protos.Response, error) {
	var buf [utils.HASHSZ]byte
	kdm_p.node.GetHash(&buf)
	var ret *protos.Response = &protos.Response{}
	ret.Hash = buf[:]

	return ret, nil
}

// DeInit calls DeInit on the internal node.
func (kdm_p *KademGo) DeInit() {
	kdm_p.node.DeInit()
}
