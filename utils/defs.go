package utils

import (
	"crypto/sha512"
	"time"
)

const (
	// KVAL is the number of K-buckets that the node will hols.
	// It also is the number of neighbours in each K-Bucket.
	KVAL = 20

	// HASHSZ defines the bytes needed to store the hash in binary, non encoded format.
	HASHSZ = sha512.Size

	// MAXHOPS is the TTL for packets that node processes. Any packet with value of 0
	// will be dropped by the node that happens to receive that.
	MAXHOPS = 20

	// PINGWAIT is the amounf of time the Ping call would wait before it fails and decides
	// to evict the neighbour in question.
	PINGWAIT = 500 * time.Millisecond
)
