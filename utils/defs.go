package utils

import "crypto/sha512"

const (
	KVAL     = 20
	HASHSZ   = sha512.Size
	GENPORT  = 12345
	PINGPORT = 12346
	MAXHOPS  = 20
)
