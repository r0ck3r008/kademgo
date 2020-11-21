package utils

import (
	"crypto/sha512"
	"net"
	"time"

	"github.com/r0ck3r008/kademgo/objmap"
	"github.com/r0ck3r008/kademgo/utils"
)

const (
	// KVAL is the number of K-buckets that the node will hols.
	// It also is the number of neighbours in each K-Bucket.
	KVAL = 20

	// ALPHAVAL is the alpha as described in the paper. This determines how many
	// peers will be called again in the recursive step of the FindPeers.
	ALPHAVAL = 3

	// HASHSZ defines the bytes needed to store the hash in binary, non encoded format.
	HASHSZ = sha512.Size

	// MAXHOPS is the TTL for packets that node processes. Any packet with value of 0
	// will be dropped by the node that happens to receive that.
	MAXHOPS = 20

	// PORTNUM defines the standard port that would be used for this application.
	PORTNUM = 12345

	// PINGWAIT is the amounf of time the Ping call would wait before it fails and decides
	// to evict the neighbour in question.
	PINGWAIT = 500 * time.Millisecond
)

// Envelope is an encapsulation which would be passed around in go channels.
type Envelope struct {
	pkt  Pkt
	addr net.UDPAddr
}

// PktType is the type that defines what the message/packet
// is supposed to represent.
type PktType int

const (
	// PingReq is the is-alive request.
	PingReq PktType = iota

	// PingRes is the response to is-alive request.
	PingRes

	// Store is the type that defines a packet to contain
	// an object that the node it supposed to store.
	Store
)

type Pkt struct {
	// Ttl is the number of hops of this particular message.
	// The message should be dropped if it ever is 0.
	Ttl int `json:"Ttl"`

	// RandNum is a long long int identifier of this message.
	// This facilitates easy caching of packets internally.
	RandNum int64 `json:"RandNum"`

	// Hash is the hash of the node sending the message.
	Hash [utils.HASHSZ]byte `json:"Hash"`

	// Type is the type of packet from PktType.
	Type PktType `json:"Type"`

	// Obj is the ObjNode type object that only Store PktType uses
	Obj objmap.ObjNode `json:"Obj"`
}
