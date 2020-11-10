package utils

import (
	"crypto/sha512"
	"io/ioutil"
)

func HashStr(data []byte) [sha512.Size]byte {
	hash := sha512.Sum512(data)

	return hash
}

func HashFile(fname *string) ([sha512.Size]byte, error) {
	data, err := ioutil.ReadFile(*fname)
	if err != nil {
		return [sha512.Size]byte{0}, err
	}

	return HashStr(data), nil
}
