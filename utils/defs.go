package utils

import (
	"crypto/sha512"
	"time"
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

	// PORTNUM defines the standard port that would be used for this application.
	PORTNUM = 12345

	// TIMEOUT is the amounf of time the Ping call would wait before it fails and decides
	// to evict the neighbour in question.
	TIMEOUT = 500 * time.Millisecond
)
