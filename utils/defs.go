package utils

import (
	"crypto/sha512"
	"time"
)

const (
	KVAL     = 20
	HASHSZ   = sha512.Size
	GENPORT  = 12345
	PINGPORT = 12346
	MAXHOPS  = 20
	PINGWAIT = 500 * time.Millisecond
)
