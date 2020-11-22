package pkt

import (
	"net"

	"github.com/r0ck3r008/kademgo/utils"
)

// ObjNode is a single bucket that stores the objects in a vector.
// Each object has its index mapped using its hash within the same node.
type ObjNode struct {
	Nmap map[[utils.HASHSZ]byte]int
	Nvec []net.IP
}

type NbrAddr struct {
	Addr net.UDPAddr
	Hash [utils.HASHSZ]byte
}

// PktType is the type that defines what the message/packet
// is supposed to represent.
type PacketType int

const (
	// PingReq is the is-alive request.
	PingReq PacketType = iota

	// PingRes is the response to is-alive request.
	PingRes

	// Store is the type that defines a packet to contain
	// an object that the node it supposed to store.
	Store
)

type Packet struct {
	// Ttl is the number of hops of this particular message.
	// The message should be dropped if it ever is 0.
	Ttl int `json:"Ttl"`

	// RandNum is a long long int identifier of this message.
	// This facilitates easy caching of packets internally.
	RandNum int64 `json:"RandNum"`

	// Hash is the hash of the node sending the message.
	Hash [utils.HASHSZ]byte `json:"Hash"`

	// Type is the type of packet from PktType.
	Type PacketType `json:"Type"`

	// Obj is the ObjNode type object that only Store PktType uses
	Obj objmap.ObjNode `json:"Obj"`
}

// Envelope is an encapsulation which would be passed around in go channels.
type Envelope struct {
	Pkt  Packet
	Addr net.UDPAddr
}
