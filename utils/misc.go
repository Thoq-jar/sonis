package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func HashString(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
