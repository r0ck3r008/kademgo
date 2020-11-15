// This package is responsible for providing general utilities that might be used by
// any of the other packages

package utils

import (
	"crypto/sha512"
	"io/ioutil"
)

func HashStr(data []byte) [HASHSZ]byte {
	hash := sha512.Sum512(data)

	return hash
}

func HashFile(fname *string) ([HASHSZ]byte, error) {
	data, err := ioutil.ReadFile(*fname)
	if err != nil {
		return [HASHSZ]byte{0}, err
	}

	return HashStr(data), nil
}
