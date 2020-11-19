// Utils package is responsible for providing general utilities that might be used by
// any of the other packages
package utils

import (
	"crypto/sha512"
	"io/ioutil"
)

// HashStr hashes the given string to SHA512 and returns the digest.
func HashStr(data []byte) [HASHSZ]byte {
	hash := sha512.Sum512(data)

	return hash
}

// HashFile hashes the given file by reading into memory and returns the digest.
func HashFile(fname *string) ([HASHSZ]byte, error) {
	data, err := ioutil.ReadFile(*fname)
	if err != nil {
		return [HASHSZ]byte{0}, err
	}

	return HashStr(data), nil
}

// GetDist is the function that calculates the `distance' of the given node's
// hash from the node which is caling it.
func GetDist(hash1 *[HASHSZ]byte, hash2 *[HASHSZ]byte) int {
	var max_pow int
	// The indx is basically the log of 2^{i} as mentioned in the algorithm
	// The algorithm states that each kbucket stores addresses with distance
	// of 2_{i} < d < 2_{i+1} where 0 <= i < 160. This indx is that `i'
	for i := HASHSZ - 1; i >= 0; i-- {
		if ((*hash1)[i] ^ (*hash2)[i]) == 1 {
			max_pow = i
			break
		}
	}

	return max_pow
}
