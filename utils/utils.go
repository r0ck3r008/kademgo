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
