// Node package is responsible for:
// 1. Creating a hash for itself
// 2. Instantiating a nbrmap object
// 3. Instantiating an objstore object
// 4. Providing the API for Kademlia RPCs
// Kademlia API for RPCs is implemented in api.go
package node

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/r0ck3r008/kademgo/connector"
	"github.com/r0ck3r008/kademgo/nbrmap"
	"github.com/r0ck3r008/kademgo/objstore"
	"github.com/r0ck3r008/kademgo/utils"
	"google.golang.org/protobuf/proto"
)

// Node structure encapsulates the UDP listening port, objstore object,
// NbrMap object as well as the hash of the node in question.
type Node struct {
	nmap *nbrmap.NbrMap
	ost  *objstore.ObjStore
	hash [utils.HASHSZ]byte
	conn *connector.Connector
}

// NodeInit is the function that initiates the ObjStore, NbrMap, UDP listener
// and as well as forms the random hash for the node.
func NodeInit(addr *string) (*Node, error) {
	node_p := &Node{}
	node_p.nmap = nbrmap.NbrMapInit()
	node_p.ost = objstore.ObjStoreInit()

	rand.Seed(time.Now().UnixNano())
	var rnum_str string = strconv.FormatInt(int64(rand.Int()), 10)
	node_p.hash = utils.HashStr([]byte(rnum_str))

	conn, err := connector.ConnectorInit(addr)
	if err != nil {
		return nil, fmt.Errorf("Error in starting listener: %s\n", err)
	}
	node_p.conn = conn

	return node_p, nil
}

// SrvLoop is the main loop that listens for messages and selects which go routine to launch
// based on the type of the message received.
func (node_p *Node) SrvLoop() {
	for {
		// Read incoming data from socket (Make sure this is the right length)
		var cmdr []byte = make([]byte, 512)
		_, addr_p, err := node_p.conn.ReadFromUDP(cmdr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in reading: %s\n", err)
			os.Exit(1)
		}

		// Unmarshal the data
		var pkt_p *Pkt = &Pkt{}
		if err := proto.Unmarshal(cmdr, pkt_p); err != nil {
			fmt.Fprintf(os.Stderr, "Error in unmarshalling: %s\n", err)
			os.Exit(1)
		}

		// Check for possible packet type
		switch pkt_p.GetType() {
		case Pkt_PingReq:
			go node_p.PingRep(addr_p)
		case Pkt_PingRep:
		default:
			fmt.Fprintf(os.Stderr, "Unknown type received!\n")
			os.Exit(1)
		}
	}
}

// DeInit closes the listening UDP connection and makes the node exit gracefully.
func (node_p *Node) DeInit() {
	node_p.conn.Close()
}
