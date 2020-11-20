package connector

import (
	"github.com/r0ck3r008/kademgo/objmap"
	"github.com/r0ck3r008/kademgo/utils"
)

// PktType is the type that defines what the message/packet
// is supposed to represent.
type PktType int

const (
	// BeginReq is the message that is sent when a node first
	// comes online
	BeginReq PktType = iota

	// BeginRes is the message that is the reply to BeginReq. Any
	// node is capable of sedning it.
	BeginRes

	// PingReq is the is-alive request.
	PingReq

	// PingRes is the response to is-alive request.
	PingRes

	// Store is the type that defines a packet to contain
	// an object that the node it supposed to store.
	Store
)

type Pkt struct {
	// Ttl is the number of hops of this particular message.
	// The message should be dropped if it ever is 0.
	Ttl int

	// RandNum is a long long int identifier of this message.
	// This facilitates easy caching of packets internally.
	RandNum int64

	// Hash is the hash of the node sending the message.
	Hash [utils.HASHSZ]byte

	// Type is the type of packet from PktType.
	Type PktType
}
